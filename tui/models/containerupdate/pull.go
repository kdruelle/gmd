package containerupdate

// func startPull(client *docker.Monitor, image string) chan commands.PullProgressMsg {

// 	ch := make(chan commands.PullProgressMsg)

// 	go func() {
// 		defer close(ch)

// 		log.Printf("pulling image %s", image)

// 		err := client.PullImageWithProgress(context.Background(), image, func(raw map[string]interface{}) {

// 			log.Printf("pull event : %+v", raw)

// 			status, _ := raw["status"].(string)
// 			layer, _ := raw["id"].(string)
// 			p, _ := raw["progress"].(string)
// 			pd, _ := raw["progressDetail"].(map[string]interface{})

// 			c, _ := pd["current"].(float64)
// 			t, _ := pd["total"].(float64)

// 			pct := 0.0
// 			if t > 0 {
// 				pct = c / t
// 			}

// 			ch <- commands.PullProgressMsg{
// 				LayerID:         layer,
// 				Status:          status,
// 				Progress:        p,
// 				ProgressCurrent: c,
// 				ProgressTotal:   t,
// 				ProgressPct:     pct,
// 			}
// 		})
// 		if err != nil {
// 			log.Printf("error pulling image %s: %v", image, err)
// 			ch <- commands.PullProgressMsg{
// 				Err: err,
// 			}
// 		}
// 	}()

// 	return ch
// }
