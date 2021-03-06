package rokuAPI

import (
	"encoding/json"
	"strings"
	"fmt"
)

func PowerEndpoint() (out map[string]string) {
	out = map[string]string{
		"success": "true",
		"roku_server_data": "None",
		"lirc_server_data": "None",
		"errors": "None",
	}

	redis_ctx := *NewRedisClient()
	roku_addr := redis_ctx.Get("roku_address")

	if roku_addr != "" {
		// first try our API endoint
		resp := RokuRequest(roku_addr, "/keypress/Power", "post")

		if strings.Contains(resp, "HTTP/1.1 200 OK") {
			out["roku_server_data"] = "Successfully hit /keypress/Power Roku API endpoint"

			return out
		} else {
			out["success"] = "false"
			out["roku_server_data"] = "Unable to hit /keypress/Power Roku API endpoint"
		}

		// if that fails (likely b/c the tv powered off and disconnected from the network), turn on using LIRC
		resp = LircServerRequest(LIRC_SERVER_ADDR, "/power", "GET")

		if strings.Contains(resp, "false") {
			out["lirc_server_data"] = "Unable to hit /power LIRC server API endpoint"
		} else {
			out["lirc_server_data"] = "Successfully hit /power LIRC server API endpoint"
			out["success"] = "true"
		}

	} else {
		out["success"] = "false"
		out["errors"] = "Unable to find Roku Address"
	}

	return out
}

func LaunchAppEndpoint(name string) (out map[string]string) {
	out = map[string]string{
		"success": "true",
		"data": "None",
	}

	fmt.Println(">>> locating id associated with app ", name)

	redis_ctx := *NewRedisClient()
	roku_addr := redis_ctx.Get("roku_address")

	if roku_addr == "" {
		out["succes"], out["data"] = "false", "Unable to locate Roku device"

		return out
	}

	// get id associated with app
	encoded_name_id_map := redis_ctx.Get("roku_apps")
	var decoded_name_id_map map[string]interface{}
	if encoded_name_id_map == "" {
		raw_apps := RokuRequest(roku_addr, "/query/apps", "GET")
		parsed_apps := ParseApps(raw_apps) // parse apps

		str, err := json.Marshal(parsed_apps)
		if err != nil {
			fmt.Println("Error encoding JSON")
			panic(err)
		}

		redis_ctx.Set("roku_apps", string(str)) // need to hash this first

	} else {
		json.Unmarshal([]byte(encoded_name_id_map), &decoded_name_id_map)
		fmt.Println(decoded_name_id_map["Netflix"])
	}

	// now we can locate the app
	var found_app_id string
	for key, value := range decoded_name_id_map {
		if strings.Contains(strings.ToLower(key), strings.ToLower(name)) {
			found_app_id = fmt.Sprintf("%s", value)
			break
		}
	}

	// launch app
	if found_app_id != "" {
		route := fmt.Sprintf("/launch/%s", found_app_id)
		RokuRequest(roku_addr, route, "POST")

	} else {
		out["success"] = "false"
		out["data"] = fmt.Sprintf("Unable to find matching app '%s'", name)
	}

	return out
}