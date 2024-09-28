package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	beego "github.com/beego/beego/v2/server/web"
)

type CatController struct {
	beego.Controller
}

type CatImage struct {
	URL string `json:"url"`
}

type Breed struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Favorite struct {
	ImageID string `json:"image_id"`
	SubID   string `json:"sub_id"`
}

type Vote struct {
	ImageID string `json:"image_id"`
	SubID   string `json:"sub_id"`
	Value   int    `json:"value"`
}

type APIResponse struct {
	Body  []byte
	Error error
}

func makeAPIRequest(method, url string, body []byte, apiKey string) <-chan APIResponse {
	responseChan := make(chan APIResponse)

	go func() {
		defer close(responseChan)

		req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
		if err != nil {
			responseChan <- APIResponse{nil, fmt.Errorf("error creating request: %v", err)}
			return
		}

		req.Header.Set("x-api-key", apiKey)
		if method == "POST" {
			req.Header.Set("Content-Type", "application/json")
		}

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			responseChan <- APIResponse{nil, fmt.Errorf("error making request: %v", err)}
			return
		}
		defer resp.Body.Close()

		responseBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			responseChan <- APIResponse{nil, fmt.Errorf("error reading response body: %v", err)}
			return
		}

		responseChan <- APIResponse{responseBody, nil}
	}()

	return responseChan
}

func (c *CatController) GetCatImages() {
	apiKey, _ := beego.AppConfig.String("cat_api_key")
	url := "https://api.thecatapi.com/v1/images/search?limit=10"

	responseChan := makeAPIRequest("GET", url, nil, apiKey)

	select {
	case response := <-responseChan:
		if response.Error != nil {
			c.Ctx.Output.SetStatus(500)
			c.Data["json"] = map[string]string{"error": response.Error.Error()}
		} else {
			var catImages []CatImage
			if err := json.Unmarshal(response.Body, &catImages); err != nil {
				c.Ctx.Output.SetStatus(500)
				c.Data["json"] = map[string]string{"error": fmt.Sprintf("Error parsing cat images: %v", err)}
			} else {
				c.Data["json"] = catImages
			}
		}
	case <-time.After(15 * time.Second):
		c.Ctx.Output.SetStatus(504)
		c.Data["json"] = map[string]string{"error": "Request timed out"}
	}

	c.ServeJSON()
}

func (c *CatController) GetBreeds() {
	url := "https://api.thecatapi.com/v1/breeds"
	responseChan := makeAPIRequest("GET", url, nil, "")

	select {
	case response := <-responseChan:
		if response.Error != nil {
			c.Ctx.Output.SetStatus(500)
			c.Data["json"] = map[string]string{"error": response.Error.Error()}
		} else {
			var breeds []Breed
			if err := json.Unmarshal(response.Body, &breeds); err != nil {
				c.Ctx.Output.SetStatus(500)
				c.Data["json"] = map[string]string{"error": fmt.Sprintf("Error parsing breeds: %v", err)}
			} else {
				c.Data["json"] = breeds
			}
		}
	case <-time.After(15 * time.Second):
		c.Ctx.Output.SetStatus(504)
		c.Data["json"] = map[string]string{"error": "Request timed out"}
	}

	c.ServeJSON()
}

func (c *CatController) GetCatImagesByBreed() {
	breedID := c.GetString("breed_id")
	if breedID == "" {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]string{"error": "Breed ID is required"}
		c.ServeJSON()
		return
	}

	url := fmt.Sprintf("https://api.thecatapi.com/v1/images/search?breed_ids=%s", breedID)
	responseChan := makeAPIRequest("GET", url, nil, "")

	select {
	case response := <-responseChan:
		if response.Error != nil {
			c.Ctx.Output.SetStatus(500)
			c.Data["json"] = map[string]string{"error": response.Error.Error()}
		} else {
			var catImages []map[string]interface{}
			if err := json.Unmarshal(response.Body, &catImages); err != nil {
				c.Ctx.Output.SetStatus(500)
				c.Data["json"] = map[string]string{"error": fmt.Sprintf("Error parsing cat images: %v", err)}
			} else {
				c.Data["json"] = catImages
			}
		}
	case <-time.After(15 * time.Second):
		c.Ctx.Output.SetStatus(504)
		c.Data["json"] = map[string]string{"error": "Request timed out"}
	}

	c.ServeJSON()
}

func (c *CatController) AddFavorite() {
	apiKey, _ := beego.AppConfig.String("cat_api_key")

	body, err := io.ReadAll(c.Ctx.Request.Body)
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]string{"error": fmt.Sprintf("Error reading request body: %v", err)}
		c.ServeJSON()
		return
	}

	var favorite Favorite
	if err := json.Unmarshal(body, &favorite); err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]string{"error": fmt.Sprintf("Error parsing request body: %v", err)}
		c.ServeJSON()
		return
	}

	url := "https://api.thecatapi.com/v1/favourites"
	responseChan := makeAPIRequest("POST", url, body, apiKey)

	select {
	case response := <-responseChan:
		if response.Error != nil {
			c.Ctx.Output.SetStatus(500)
			c.Data["json"] = map[string]string{"error": response.Error.Error()}
		} else {
			c.Ctx.Output.SetStatus(200)
			c.Ctx.Output.Body(response.Body)
		}
	case <-time.After(15 * time.Second):
		c.Ctx.Output.SetStatus(504)
		c.Data["json"] = map[string]string{"error": "Request timed out"}
	}

	c.ServeJSON()
}

func (c *CatController) GetFavorites() {
	apiKey, _ := beego.AppConfig.String("cat_api_key")
	subID := c.GetString("sub_id")
	url := fmt.Sprintf("https://api.thecatapi.com/v1/favourites?sub_id=%s", subID)

	responseChan := makeAPIRequest("GET", url, nil, apiKey)

	select {
	case response := <-responseChan:
		if response.Error != nil {
			c.Ctx.Output.SetStatus(500)
			c.Data["json"] = map[string]string{"error": response.Error.Error()}
		} else {
			c.Ctx.Output.SetStatus(200)
			c.Ctx.Output.Body(response.Body)
		}
	case <-time.After(15 * time.Second):
		c.Ctx.Output.SetStatus(504)
		c.Data["json"] = map[string]string{"error": "Request timed out"}
	}

	c.ServeJSON()
}

func (c *CatController) DeleteFavorite() {
	apiKey, _ := beego.AppConfig.String("cat_api_key")
	favoriteID := c.Ctx.Input.Param(":id")
	url := fmt.Sprintf("https://api.thecatapi.com/v1/favourites/%s", favoriteID)

	responseChan := makeAPIRequest("DELETE", url, nil, apiKey)

	select {
	case response := <-responseChan:
		if response.Error != nil {
			c.Ctx.Output.SetStatus(500)
			c.Data["json"] = map[string]string{"error": response.Error.Error()}
		} else {
			c.Ctx.Output.SetStatus(200)
			c.Ctx.Output.Body(response.Body)
		}
	case <-time.After(15 * time.Second):
		c.Ctx.Output.SetStatus(504)
		c.Data["json"] = map[string]string{"error": "Request timed out"}
	}

	c.ServeJSON()
}

func (c *CatController) Vote() {
	apiKey, _ := beego.AppConfig.String("cat_api_key")

	body, err := io.ReadAll(c.Ctx.Request.Body)
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]string{"error": fmt.Sprintf("Error reading request body: %v", err)}
		c.ServeJSON()
		return
	}

	var vote Vote
	if err := json.Unmarshal(body, &vote); err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]string{"error": fmt.Sprintf("Error parsing request body: %v", err)}
		c.ServeJSON()
		return
	}

	url := "https://api.thecatapi.com/v1/votes"
	responseChan := makeAPIRequest("POST", url, body, apiKey)

	select {
	case response := <-responseChan:
		if response.Error != nil {
			c.Ctx.Output.SetStatus(500)
			c.Data["json"] = map[string]string{"error": response.Error.Error()}
		} else {
			c.Ctx.Output.SetStatus(200)
			c.Ctx.Output.Body(response.Body)
		}
	case <-time.After(15 * time.Second):
		c.Ctx.Output.SetStatus(504)
		c.Data["json"] = map[string]string{"error": "Request timed out"}
	}

	c.ServeJSON()
}

func (c *CatController) GetVotes() {
	apiKey, _ := beego.AppConfig.String("cat_api_key")
	limit := c.GetString("limit")
	order := c.GetString("order")
	subID := c.GetString("sub_id")
	page := c.GetString("page")

	url := fmt.Sprintf("https://api.thecatapi.com/v1/votes?limit=%s&order=%s&sub_id=%s&page=%s", limit, order, subID, page)

	responseChan := makeAPIRequest("GET", url, nil, apiKey)

	select {
	case response := <-responseChan:
		if response.Error != nil {
			c.Ctx.Output.SetStatus(500)
			c.Data["json"] = map[string]string{"error": response.Error.Error()}
		} else {
			c.Ctx.Output.SetStatus(200)
			c.Ctx.Output.Body(response.Body)
		}
	case <-time.After(15 * time.Second):
		c.Ctx.Output.SetStatus(504)
		c.Data["json"] = map[string]string{"error": "Request timed out"}
	}

	c.ServeJSON()
}
