/*
 * Testing
 *
 * No description provided (generated by Swagger Codegen https://github.com/swagger-api/swagger-codegen)
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

type UsersAdded struct {
	Status string `json:"status"`
	NumRequested int32 `json:"numRequested"`
	NumExisted int32 `json:"numExisted"`
	NumAdded int32 `json:"numAdded"`
}
