package mongodbatlas

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/go-test/deep"
	"github.com/openlyinc/pointy"
)

func TestDataLakes_List(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	groupID := "6c7498dg87d9e6526801572b"

	path := fmt.Sprintf("/groups/%s/dataLakes", groupID)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `[
			{
				"cloudProviderConfig": {
					"aws": {
					  "iamAssumedRoleARN": "arn:aws:iam::123456789012:role/ReadS3BucketRole"
				  	}
			  	},
			  	"dataProcessRegion": {
					"cloudProvider" : "AWS",
					"region" : "VIRGINIA_USA"
			  	},
			  	"groupId": "6c7498dg87d9e6526801572b",
			  	"hostnames": [
					"usermetricdata.mongodb.example.net"
			  	],
			  	"name": "UserMetricData",
			  	"state": "ACTIVE",
			  	"storage": {
				  	"databases": [
						{
							"name": "my.database",
							"collections": [
								{
									"name": "my.collection",
									"dataSources": [
										{
												"storeName" : "store",
												"defaultFormat" : ".json",
												"path" : "/path"
										}
									]
								}
							],
							"views": [
								{
									"name" : "my.view",
									"source" : "source",
									"pipeline" : "my.pipeline"
								}
							]
						}
					],
					"stores": [
						{
							"name": "datacenter-alpha",
							"provider": "s3",
						  	"region": "us-east-1",
						  	"bucket": "datacenter-alpha",
						  	"prefix": "/metrics",
						  	"delimiter": "/",
						  	"includeTags": false
						}
					]
				}
			}
		]`)
	})

	snapshots, _, err := client.DataLakes.List(ctx, groupID)
	if err != nil {
		t.Fatalf("DataLake.List returned error: %v", err)
	}

	expected := []DataLake{
		{
			CloudProviderConfig: CloudProviderConfig{
				AWSConfig: AwsCloudProviderConfig{
					IAMAssumedRoleARN: "arn:aws:iam::123456789012:role/ReadS3BucketRole",
					TestS3Bucket:      "",
				},
			},
			DataProcessRegion: DataProcessRegion{
				CloudProvider: "AWS",
				Region:        "VIRGINIA_USA",
			},
			GroupID:   groupID,
			Hostnames: []string{"usermetricdata.mongodb.example.net"},
			Name:      "UserMetricData",
			State:     "ACTIVE",
			Storage: Storage{
				Databases: []DataLakeDatabase{
					{
						Name: "my.database",
						Collections: []DataLakeCollection{
							{
								Name: "my.collection",
								DataSources: []DataLakeDataSource{
									{
										StoreName:     "store",
										DefaultFormat: ".json",
										Path:          "/path",
									},
								},
							},
						},
						Views: []DataLakeDatabaseView{
							{
								Name:     "my.view",
								Source:   "source",
								Pipeline: "my.pipeline",
							},
						},
					},
				},
				Stores: []DataLakeStore{
					{
						Name:        "datacenter-alpha",
						Provider:    "s3",
						Region:      "us-east-1",
						Bucket:      "datacenter-alpha",
						Prefix:      "/metrics",
						Delimiter:   "/",
						IncludeTags: pointy.Bool(false),
					},
				},
			},
		},
	}

	if diff := deep.Equal(snapshots, expected); diff != nil {
		t.Error(diff)
	}
}

func TestDataLake_Get(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	groupID := "6c7498dg87d9e6526801572b"
	dataLakeName := "UserMetricData"
	path := fmt.Sprintf("/groups/%s/dataLakes/%s", groupID, dataLakeName)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
			  "cloudProviderConfig": {
				  "aws": {
					  "iamAssumedRoleARN": "arn:aws:iam::123456789012:role/ReadS3BucketRole"
				  }
			  },
			  "dataProcessRegion": {
				"cloudProvider" : "AWS",
				"region" : "VIRGINIA_USA"
			  },
			  "groupId": "6c7498dg87d9e6526801572b",
			  "hostnames": [
				  "usermetricdata.mongodb.example.net"
			  ],
			  "name": "UserMetricData",
			  "state": "ACTIVE",
			  "storage": {
				  "databases": [],
				  "stores": []
			  }
		}`)
	})

	cloudProviderSnapshot, _, err := client.DataLakes.Get(ctx, groupID, dataLakeName)
	if err != nil {
		t.Fatalf("DataLake.Get returned error: %v", err)
	}

	expected := DataLake{
		CloudProviderConfig: CloudProviderConfig{
			AWSConfig: AwsCloudProviderConfig{
				IAMAssumedRoleARN: "arn:aws:iam::123456789012:role/ReadS3BucketRole",
				TestS3Bucket:      "",
			},
		},
		DataProcessRegion: DataProcessRegion{
			CloudProvider: "AWS",
			Region:        "VIRGINIA_USA",
		},
		GroupID:   groupID,
		Hostnames: []string{"usermetricdata.mongodb.example.net"},
		Name:      "UserMetricData",
		State:     "ACTIVE",
		Storage: Storage{
			Databases: []DataLakeDatabase{},
			Stores:    []DataLakeStore{},
		},
	}

	if diff := deep.Equal(cloudProviderSnapshot, &expected); diff != nil {
		t.Error(diff)
	}
}

func TestDataLake_Update(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	groupID := "6c7498dg87d9e6526801572b"
	dataLakeName := "UserMetricData"

	updateRequest := &DataLakeUpdateRequest{
		CloudProviderConfig: &CloudProviderConfig{
			AWSConfig: AwsCloudProviderConfig{
				IAMAssumedRoleARN: "new_arn",
				TestS3Bucket:      "new_bucket",
			},
		},
		DataProcessRegion: &DataProcessRegion{
			CloudProvider: "AWS",
			Region:        "DUBLIN_IRL",
		},
	}

	path := fmt.Sprintf("/groups/%s/dataLakes/%s", groupID, dataLakeName)
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		expected := map[string]interface{}{
			"cloudProviderConfig": map[string]interface{}{
				"aws": map[string]interface{}{
					"iamAssumedRoleARN": "new_arn",
					"testS3Bucket":      "new_bucket",
				},
			},
			"dataProcessRegion": map[string]interface{}{
				"cloudProvider": "AWS",
				"region":        "DUBLIN_IRL",
			},
		}

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("Decode json: %v", err)
		}

		if diff := deep.Equal(v, expected); diff != nil {
			t.Error(diff)
		}

		fmt.Fprint(w, `{
			  "cloudProviderConfig": {
				  "aws": {
					  "iamAssumedRoleARN": "new_arn",
					  "testS3Bucket": "new_bucket"
				  }
			  },
			  "dataProcessRegion": {
				"cloudProvider" : "AWS",
				"region" : "DUBLIN_IRL"
			  },
			  "groupId": "6c7498dg87d9e6526801572b",
			  "hostnames": [
				  "usermetricdata.mongodb.example.net"
			  ],
			  "name": "UserMetricData",
			  "state": "ACTIVE",
			  "storage": {
				  "databases": [],
				  "stores": []
			  }
		}`)
	})

	updatedDataLake, _, err := client.DataLakes.Update(ctx, groupID, dataLakeName, updateRequest)
	if err != nil {
		t.Fatalf("DataLake.Update returned error: %v", err)
	}

	expected := DataLake{
		CloudProviderConfig: CloudProviderConfig{
			AWSConfig: AwsCloudProviderConfig{
				IAMAssumedRoleARN: "new_arn",
				TestS3Bucket:      "new_bucket",
			},
		},
		DataProcessRegion: DataProcessRegion{
			CloudProvider: "AWS",
			Region:        "DUBLIN_IRL",
		},
		GroupID:   groupID,
		Hostnames: []string{"usermetricdata.mongodb.example.net"},
		Name:      "UserMetricData",
		State:     "ACTIVE",
		Storage: Storage{
			Databases: []DataLakeDatabase{},
			Stores:    []DataLakeStore{},
		},
	}

	if diff := deep.Equal(updatedDataLake, &expected); diff != nil {
		t.Error(diff)
	}
}

func TestDataLake_Create(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	groupID := "6c7498dg87d9e6526801572b"
	dataLakeName := "UserMetricData"

	createRequest := &DataLakeCreateRequest{
		Name: dataLakeName,
	}

	path := fmt.Sprintf("/groups/%s/dataLakes", groupID)
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		expected := map[string]interface{}{
			"name": "UserMetricData",
		}

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("Decode json: %v", err)
		}

		if diff := deep.Equal(v, expected); diff != nil {
			t.Error(diff)
		}

		fmt.Fprint(w, `{
			  "cloudProviderConfig": {
				  "aws": {
					  "iamAssumedRoleARN": "new_arn",
					  "testS3Bucket": "new_bucket"
				  }
			  },
			  "dataProcessRegion": {
				"cloudProvider" : "AWS",
				"region" : "DUBLIN_IRL"
			  },
			  "groupId": "6c7498dg87d9e6526801572b",
			  "hostnames": [
				  "usermetricdata.mongodb.example.net"
			  ],
			  "name": "UserMetricData",
			  "state": "UNVERIFIED",
			  "storage": {
				  "databases": [],
				  "stores": []
			  }
		}`)
	})

	updatedDataLake, _, err := client.DataLakes.Create(ctx, groupID, createRequest)
	if err != nil {
		t.Fatalf("DataLake.Create returned error: %v", err)
	}

	expected := DataLake{
		CloudProviderConfig: CloudProviderConfig{
			AWSConfig: AwsCloudProviderConfig{
				IAMAssumedRoleARN: "new_arn",
				TestS3Bucket:      "new_bucket",
			},
		},
		DataProcessRegion: DataProcessRegion{
			CloudProvider: "AWS",
			Region:        "DUBLIN_IRL",
		},
		GroupID:   groupID,
		Hostnames: []string{"usermetricdata.mongodb.example.net"},
		Name:      "UserMetricData",
		State:     "UNVERIFIED",
		Storage: Storage{
			Databases: []DataLakeDatabase{},
			Stores:    []DataLakeStore{},
		},
	}

	if diff := deep.Equal(updatedDataLake, &expected); diff != nil {
		t.Error(diff)
	}
}

func TestDataLake_Delete(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	groupID := "6c7498dg87d9e6526801572b"
	dataLakeName := "dataLake"

	path := fmt.Sprintf("/groups/%s/dataLakes/%s", groupID, dataLakeName)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	_, err := client.DataLakes.Delete(ctx, groupID, dataLakeName)
	if err != nil {
		t.Fatalf("DataLakes.Delete returned error: %v", err)
	}
}