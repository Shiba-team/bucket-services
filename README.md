# bucket-services
![#f03c15](https://via.placeholder.com/15/f03c15/000000?text=+) `Note`
```diff
Always put Header for all request
+ Authorizaion: Bearer + Token
```
---------
# List Function
+ [Create Bucket](https://github.com/Shiba-team/bucket-services/new/master?readme=1#create-bucket)
+ [Get List Bucket Of User](https://github.com/Shiba-team/bucket-services/new/master?readme=1#create-bucket)
+ [Get Bucket Infomation]()
+ [Add File To Bucket](https://github.com/Shiba-team/bucket-services/new/master?readme=1#create-bucket)
+ [Update Bucket Permission And Status]()
+ [Delete Bucket]()
+ [Get Bucket's Size]()
+ [Get Files In Bucket]()
+ [Get File Information]()
+ [Download File]()
+ [Update File Permission]()
+ [Update File Status]()
+ [Delete File]()
---------
## Create Bucket
### Method
`POST /4`
### Status Code
Success: `200`<br>
Fail: `400`
### Example Body 
```json
{
    "bucketname": "abcdej"
}
```
### Example Success Response
```json
{
    "_id": "60ad361df2f1f4c96afce4ea",
    "bucketname": "abcdej"
}
```
### Example for Failure Response
```json
{
    "error": "bucket name is used"
}
```
----
## Get User's Buckets
### Method
`GET /4`
### Status Code
Success: `200`<br>
Fail: `400`
### Example Success Response
```json
{
    "count": 1,
    "result": [
        {
            "_id": "60ad361df2f1f4c96afce4ea",
            "owner": "admin",
            "bucketname": "abcdej",
            "createdat": "2021-05-25T17:38:37.858Z",
            "listfile": null,
            "status": 0,
            "permission": 0,
            "lastmodified": "2021-05-25T17:38:37.858Z"
        }
    ]
}
```
------------------------------------------------------------------------------------------
## Add File To Bucket
### Method
`POST /4/:BucketID`
### Status Code
Success: `200`<br>
Fail: `400`
### Example Body 
```json
{
    "file": "<binaryfile>",
    "filename": "name.pdf"
}
```
### Example Success Response
filepath is download link 
```json
{
    "filename": "01c42d0d51d9d62955d7aa507bd75da38a9666f7decbe05fbc4c7990a157d610a82642d91bf25fc607e021a80f7e9729285e5ddb234beb2cfabad91d6de2c6fc.pdf",
    "filepath": "http://localhost:3000/4/60ad361df2f1f4c96afce4ea/0/01c42d0d51d9d62955d7aa507bd75da38a9666f7decbe05fbc4c7990a157d610a82642d91bf25fc607e021a80f7e9729285e5ddb234beb2cfabad91d6de2c6fc.pdf"
}
```
### Example for Failure Response
If you don't choose file, server will be response
```json
{
    "error": "http: no such file"
}
```
-------------------------------------------------------------------
## Get Bucket Information
### Method
`GET /4/:BucketID`
### Status Code
Success: `200`<br>
Fail: `400`
### Example Success Response
```json
{
    "_id": "60ad361df2f1f4c96afce4ea",
    "owner": "admin",
    "bucketname": "abcdej",
    "createdat": "2021-05-25T17:38:37.858Z",
    "listfile": null,
    "status": 0,
    "permission": 0,
    "lastmodified": "2021-05-25T17:38:37.858Z"
}
```
### Example for Failure Response
If bucketid is wrongs, server will be response
```json
{
    "error": "encoding/hex: odd length hex string"
}
```
-----------------------------------------------------------------------
## Update Bucket Permission And Status
if wrongs value, bucket information do not change
### Method
`PATCH /4/:BucketID`
### Status Code
Success: `200`
### Example Body 
- Permission:
  + 0: Read
  + 1: Write
  + 2: ReadAndWrite
- Status:
  + 0: Public
  + 1: Private

```json
{
    "permission":2,
    "status":1
}
```
### Example Success Response
```json
{
    "message": "Updated"
}
```
------------------------------------------------------------------------------------------
## Delete Bucket
### Method
`Delete /4/:BucketID`
### Status Code
Success: `200`
### Example Success Response
```json
{
    "message": "abcdej deleted successfully"
}
```
------------------------------------------------------------------------------------------
## Get Size Of Bucket
### Method
`GET /4/:BucketID/1`
### Status Code
Success: `200`<br>
### Example Success Response
```json
{
    "bucketsize": 5737
}
```
------------------------------------------------------------------------------------------
## Get List File Of Bucket
### Method
`Get /4/:BucketID/2`
### Status Code
Success: `200`
### Example Success Response
```json
{
    "count": 1,
    "result": [
        {
            "filelink": "",
            "downloadlink": "http://localhost:3000/4/6096e5124b24030d230ab107/0/eb535076378fd96a1db0f8fb20f3830ae6fc624f02e2d673958a433f1ffaf78d1208370f274c66d7f6b905f1cff9d7a4440d9c0ff6e6eb2892ead6544a058c00.png",
            "createdby": "testuser",
            "filename": "test.png",
            "s3name": "eb535076378fd96a1db0f8fb20f3830ae6fc624f02e2d673958a433f1ffaf78d1208370f274c66d7f6b905f1cff9d7a4440d9c0ff6e6eb2892ead6544a058c00.png",
            "createdat": "2021-05-08T19:24:40.528Z",
            "lastmodify": "2021-05-08T19:24:40.528Z",
            "status": 1,
            "permission": 2,
            "size": 5737
        }
    ]
}
------------------------------------------------------------------------------------------
## Get File Info
### Method
`GET /4/:BucketID/1/:filename`
### Status Code
Success: `200`
### Example Success Response
`s3name` is filename in `URL`
```json
{
    "filelink": "",
    "downloadlink": "http://localhost:3000/4/6096e5124b24030d230ab107/0/eb535076378fd96a1db0f8fb20f3830ae6fc624f02e2d673958a433f1ffaf78d1208370f274c66d7f6b905f1cff9d7a4440d9c0ff6e6eb2892ead6544a058c00.png",
    "createdby": "testuser",
    "filename": "test.png",
    "s3name": "eb535076378fd96a1db0f8fb20f3830ae6fc624f02e2d673958a433f1ffaf78d1208370f274c66d7f6b905f1cff9d7a4440d9c0ff6e6eb2892ead6544a058c00.png",
    "createdat": "2021-05-08T19:24:40.528Z",
    "lastmodify": "2021-05-08T19:24:40.528Z",
    "status": 1,
    "permission": 2,
    "size": 5737
}
------------------------------------------------------------------------------------------
## Download File In Bucket
Download file needn't `Authorization` if it is `directlink` (Status = 0 and Permission = 0)
### Method
`GET /4/:BucketID/0/:filename`
### Status Code
Success: `200`
### Example Success Response
`Binary File`
------------------------------------------------------------------------------------------
## Change File Permission
### Method
`PATCH /4/:BucketID/2/:filename`
### Status Code
Success: `200`
### Example Body 
see [permission](https://github.com/Shiba-team/bucket-services/new/master?readme=1#example-body-2)
```json
{
    "permission":2
}
```
### Example Success Response
```json
{
    "message": "permission updated successfully"
}
```
------------------------------------------------------------------------------------------
## Change File Status
### Method
`PATCH /4/:BucketID/3/:filename`
### Status Code
Success: `200`
### Example Body 
see [status](https://github.com/Shiba-team/bucket-services/new/master?readme=1#example-body-2)
```json
{
    "status":1
}
```
### Example Success Response
```json
{
    "message": "status updated successfully"
}
```
------------------------------------------------------------------------------------------
## Delete File
### Method
`POST /4/:BucketID/4/:filename`
### Status Code
Success: `200`
### Example Success Response
```json
{
    "message": "bfddcbe16a788e50c50f27436ee782ba920f8d5ff5d611b9c40310afe92d363364a650fb16cafdc70a377034e1332afbb07d7bf23a8572e3194d922d7cea3b28.txt has removed"
}
```
-----------------------------------------------------------------------------
