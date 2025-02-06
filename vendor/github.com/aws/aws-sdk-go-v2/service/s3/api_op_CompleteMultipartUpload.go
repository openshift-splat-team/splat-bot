// Code generated by smithy-go-codegen DO NOT EDIT.

package s3

import (
	"context"
	"fmt"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	s3cust "github.com/aws/aws-sdk-go-v2/service/s3/internal/customizations"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Completes a multipart upload by assembling previously uploaded parts.
//
// You first initiate the multipart upload and then upload all parts using the [UploadPart]
// operation or the [UploadPartCopy]operation. After successfully uploading all relevant parts of
// an upload, you call this CompleteMultipartUpload operation to complete the
// upload. Upon receiving this request, Amazon S3 concatenates all the parts in
// ascending order by part number to create a new object. In the
// CompleteMultipartUpload request, you must provide the parts list and ensure that
// the parts list is complete. The CompleteMultipartUpload API operation
// concatenates the parts that you provide in the list. For each part in the list,
// you must provide the PartNumber value and the ETag value that are returned
// after that part was uploaded.
//
// The processing of a CompleteMultipartUpload request could take several minutes
// to finalize. After Amazon S3 begins processing the request, it sends an HTTP
// response header that specifies a 200 OK response. While processing is in
// progress, Amazon S3 periodically sends white space characters to keep the
// connection from timing out. A request could fail after the initial 200 OK
// response has been sent. This means that a 200 OK response can contain either a
// success or an error. The error response might be embedded in the 200 OK
// response. If you call this API operation directly, make sure to design your
// application to parse the contents of the response and handle it appropriately.
// If you use Amazon Web Services SDKs, SDKs handle this condition. The SDKs detect
// the embedded error and apply error handling per your configuration settings
// (including automatically retrying the request as appropriate). If the condition
// persists, the SDKs throw an exception (or, for the SDKs that don't use
// exceptions, they return an error).
//
// Note that if CompleteMultipartUpload fails, applications should be prepared to
// retry any failed requests (including 500 error responses). For more information,
// see [Amazon S3 Error Best Practices].
//
// You can't use Content-Type: application/x-www-form-urlencoded for the
// CompleteMultipartUpload requests. Also, if you don't provide a Content-Type
// header, CompleteMultipartUpload can still return a 200 OK response.
//
// For more information about multipart uploads, see [Uploading Objects Using Multipart Upload] in the Amazon S3 User Guide.
//
// Directory buckets - For directory buckets, you must make requests for this API
// operation to the Zonal endpoint. These endpoints support virtual-hosted-style
// requests in the format
// https://bucket_name.s3express-az_id.region.amazonaws.com/key-name . Path-style
// requests are not supported. For more information, see [Regional and Zonal endpoints]in the Amazon S3 User
// Guide.
//
// Permissions
//   - General purpose bucket permissions - For information about permissions
//     required to use the multipart upload API, see [Multipart Upload and Permissions]in the Amazon S3 User Guide.
//
// If you provide an [additional checksum value]in your MultipartUpload requests and the object is encrypted
//
//	with Key Management Service, you must have permission to use the kms:Decrypt
//	action for the CompleteMultipartUpload request to succeed.
//
//	- Directory bucket permissions - To grant access to this API operation on a
//	directory bucket, we recommend that you use the [CreateSession]CreateSession API operation
//	for session-based authorization. Specifically, you grant the
//	s3express:CreateSession permission to the directory bucket in a bucket policy
//	or an IAM identity-based policy. Then, you make the CreateSession API call on
//	the bucket to obtain a session token. With the session token in your request
//	header, you can make API requests to this operation. After the session token
//	expires, you make another CreateSession API call to generate a new session
//	token for use. Amazon Web Services CLI or SDKs create session and refresh the
//	session token automatically to avoid service interruptions when a session
//	expires. For more information about authorization, see [CreateSession]CreateSession .
//
// If the object is encrypted with SSE-KMS, you must also have the
//
//	kms:GenerateDataKey and kms:Decrypt permissions in IAM identity-based policies
//	and KMS key policies for the KMS key.
//
// Special errors
//
//   - Error Code: EntityTooSmall
//
//   - Description: Your proposed upload is smaller than the minimum allowed
//     object size. Each part must be at least 5 MB in size, except the last part.
//
//   - HTTP Status Code: 400 Bad Request
//
//   - Error Code: InvalidPart
//
//   - Description: One or more of the specified parts could not be found. The
//     part might not have been uploaded, or the specified ETag might not have matched
//     the uploaded part's ETag.
//
//   - HTTP Status Code: 400 Bad Request
//
//   - Error Code: InvalidPartOrder
//
//   - Description: The list of parts was not in ascending order. The parts list
//     must be specified in order by part number.
//
//   - HTTP Status Code: 400 Bad Request
//
//   - Error Code: NoSuchUpload
//
//   - Description: The specified multipart upload does not exist. The upload ID
//     might be invalid, or the multipart upload might have been aborted or completed.
//
//   - HTTP Status Code: 404 Not Found
//
// HTTP Host header syntax  Directory buckets - The HTTP Host header syntax is
// Bucket_name.s3express-az_id.region.amazonaws.com .
//
// The following operations are related to CompleteMultipartUpload :
//
// [CreateMultipartUpload]
//
// [UploadPart]
//
// [AbortMultipartUpload]
//
// [ListParts]
//
// [ListMultipartUploads]
//
// [Uploading Objects Using Multipart Upload]: https://docs.aws.amazon.com/AmazonS3/latest/dev/uploadobjusingmpu.html
// [Amazon S3 Error Best Practices]: https://docs.aws.amazon.com/AmazonS3/latest/dev/ErrorBestPractices.html
// [AbortMultipartUpload]: https://docs.aws.amazon.com/AmazonS3/latest/API/API_AbortMultipartUpload.html
// [ListParts]: https://docs.aws.amazon.com/AmazonS3/latest/API/API_ListParts.html
// [UploadPart]: https://docs.aws.amazon.com/AmazonS3/latest/API/API_UploadPart.html
// [Regional and Zonal endpoints]: https://docs.aws.amazon.com/AmazonS3/latest/userguide/s3-express-Regions-and-Zones.html
// [additional checksum value]: https://docs.aws.amazon.com/AmazonS3/latest/API/API_Checksum.html
// [ListMultipartUploads]: https://docs.aws.amazon.com/AmazonS3/latest/API/API_ListMultipartUploads.html
// [UploadPartCopy]: https://docs.aws.amazon.com/AmazonS3/latest/API/API_UploadPartCopy.html
// [Multipart Upload and Permissions]: https://docs.aws.amazon.com/AmazonS3/latest/dev/mpuAndPermissions.html
// [CreateMultipartUpload]: https://docs.aws.amazon.com/AmazonS3/latest/API/API_CreateMultipartUpload.html
//
// [CreateSession]: https://docs.aws.amazon.com/AmazonS3/latest/API/API_CreateSession.html
func (c *Client) CompleteMultipartUpload(ctx context.Context, params *CompleteMultipartUploadInput, optFns ...func(*Options)) (*CompleteMultipartUploadOutput, error) {
	if params == nil {
		params = &CompleteMultipartUploadInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "CompleteMultipartUpload", params, optFns, c.addOperationCompleteMultipartUploadMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*CompleteMultipartUploadOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type CompleteMultipartUploadInput struct {

	// Name of the bucket to which the multipart upload was initiated.
	//
	// Directory buckets - When you use this operation with a directory bucket, you
	// must use virtual-hosted-style requests in the format
	// Bucket_name.s3express-az_id.region.amazonaws.com . Path-style requests are not
	// supported. Directory bucket names must be unique in the chosen Availability
	// Zone. Bucket names must follow the format bucket_base_name--az-id--x-s3 (for
	// example, DOC-EXAMPLE-BUCKET--usw2-az1--x-s3 ). For information about bucket
	// naming restrictions, see [Directory bucket naming rules]in the Amazon S3 User Guide.
	//
	// Access points - When you use this action with an access point, you must provide
	// the alias of the access point in place of the bucket name or specify the access
	// point ARN. When using the access point ARN, you must direct requests to the
	// access point hostname. The access point hostname takes the form
	// AccessPointName-AccountId.s3-accesspoint.Region.amazonaws.com. When using this
	// action with an access point through the Amazon Web Services SDKs, you provide
	// the access point ARN in place of the bucket name. For more information about
	// access point ARNs, see [Using access points]in the Amazon S3 User Guide.
	//
	// Access points and Object Lambda access points are not supported by directory
	// buckets.
	//
	// S3 on Outposts - When you use this action with Amazon S3 on Outposts, you must
	// direct requests to the S3 on Outposts hostname. The S3 on Outposts hostname
	// takes the form
	// AccessPointName-AccountId.outpostID.s3-outposts.Region.amazonaws.com . When you
	// use this action with S3 on Outposts through the Amazon Web Services SDKs, you
	// provide the Outposts access point ARN in place of the bucket name. For more
	// information about S3 on Outposts ARNs, see [What is S3 on Outposts?]in the Amazon S3 User Guide.
	//
	// [Directory bucket naming rules]: https://docs.aws.amazon.com/AmazonS3/latest/userguide/directory-bucket-naming-rules.html
	// [What is S3 on Outposts?]: https://docs.aws.amazon.com/AmazonS3/latest/userguide/S3onOutposts.html
	// [Using access points]: https://docs.aws.amazon.com/AmazonS3/latest/userguide/using-access-points.html
	//
	// This member is required.
	Bucket *string

	// Object key for which the multipart upload was initiated.
	//
	// This member is required.
	Key *string

	// ID for the initiated multipart upload.
	//
	// This member is required.
	UploadId *string

	// This header can be used as a data integrity check to verify that the data
	// received is the same data that was originally sent. This header specifies the
	// base64-encoded, 32-bit CRC-32 checksum of the object. For more information, see [Checking object integrity]
	// in the Amazon S3 User Guide.
	//
	// [Checking object integrity]: https://docs.aws.amazon.com/AmazonS3/latest/userguide/checking-object-integrity.html
	ChecksumCRC32 *string

	// This header can be used as a data integrity check to verify that the data
	// received is the same data that was originally sent. This header specifies the
	// base64-encoded, 32-bit CRC-32C checksum of the object. For more information, see
	// [Checking object integrity]in the Amazon S3 User Guide.
	//
	// [Checking object integrity]: https://docs.aws.amazon.com/AmazonS3/latest/userguide/checking-object-integrity.html
	ChecksumCRC32C *string

	// This header can be used as a data integrity check to verify that the data
	// received is the same data that was originally sent. This header specifies the
	// base64-encoded, 160-bit SHA-1 digest of the object. For more information, see [Checking object integrity]
	// in the Amazon S3 User Guide.
	//
	// [Checking object integrity]: https://docs.aws.amazon.com/AmazonS3/latest/userguide/checking-object-integrity.html
	ChecksumSHA1 *string

	// This header can be used as a data integrity check to verify that the data
	// received is the same data that was originally sent. This header specifies the
	// base64-encoded, 256-bit SHA-256 digest of the object. For more information, see [Checking object integrity]
	// in the Amazon S3 User Guide.
	//
	// [Checking object integrity]: https://docs.aws.amazon.com/AmazonS3/latest/userguide/checking-object-integrity.html
	ChecksumSHA256 *string

	// The account ID of the expected bucket owner. If the account ID that you provide
	// does not match the actual owner of the bucket, the request fails with the HTTP
	// status code 403 Forbidden (access denied).
	ExpectedBucketOwner *string

	// Uploads the object only if the object key name does not already exist in the
	// bucket specified. Otherwise, Amazon S3 returns a 412 Precondition Failed error.
	//
	// If a conflicting operation occurs during the upload S3 returns a 409
	// ConditionalRequestConflict response. On a 409 failure you should re-initiate the
	// multipart upload with CreateMultipartUpload and re-upload each part.
	//
	// Expects the '*' (asterisk) character.
	//
	// For more information about conditional requests, see [RFC 7232], or [Conditional requests] in the Amazon S3
	// User Guide.
	//
	// [Conditional requests]: https://docs.aws.amazon.com/AmazonS3/latest/userguide/conditional-requests.html
	// [RFC 7232]: https://tools.ietf.org/html/rfc7232
	IfNoneMatch *string

	// The container for the multipart upload request information.
	MultipartUpload *types.CompletedMultipartUpload

	// Confirms that the requester knows that they will be charged for the request.
	// Bucket owners need not specify this parameter in their requests. If either the
	// source or destination S3 bucket has Requester Pays enabled, the requester will
	// pay for corresponding charges to copy the object. For information about
	// downloading objects from Requester Pays buckets, see [Downloading Objects in Requester Pays Buckets]in the Amazon S3 User
	// Guide.
	//
	// This functionality is not supported for directory buckets.
	//
	// [Downloading Objects in Requester Pays Buckets]: https://docs.aws.amazon.com/AmazonS3/latest/dev/ObjectsinRequesterPaysBuckets.html
	RequestPayer types.RequestPayer

	// The server-side encryption (SSE) algorithm used to encrypt the object. This
	// parameter is required only when the object was created using a checksum
	// algorithm or if your bucket policy requires the use of SSE-C. For more
	// information, see [Protecting data using SSE-C keys]in the Amazon S3 User Guide.
	//
	// This functionality is not supported for directory buckets.
	//
	// [Protecting data using SSE-C keys]: https://docs.aws.amazon.com/AmazonS3/latest/userguide/ServerSideEncryptionCustomerKeys.html#ssec-require-condition-key
	SSECustomerAlgorithm *string

	// The server-side encryption (SSE) customer managed key. This parameter is needed
	// only when the object was created using a checksum algorithm. For more
	// information, see [Protecting data using SSE-C keys]in the Amazon S3 User Guide.
	//
	// This functionality is not supported for directory buckets.
	//
	// [Protecting data using SSE-C keys]: https://docs.aws.amazon.com/AmazonS3/latest/dev/ServerSideEncryptionCustomerKeys.html
	SSECustomerKey *string

	// The MD5 server-side encryption (SSE) customer managed key. This parameter is
	// needed only when the object was created using a checksum algorithm. For more
	// information, see [Protecting data using SSE-C keys]in the Amazon S3 User Guide.
	//
	// This functionality is not supported for directory buckets.
	//
	// [Protecting data using SSE-C keys]: https://docs.aws.amazon.com/AmazonS3/latest/dev/ServerSideEncryptionCustomerKeys.html
	SSECustomerKeyMD5 *string

	noSmithyDocumentSerde
}

func (in *CompleteMultipartUploadInput) bindEndpointParams(p *EndpointParameters) {

	p.Bucket = in.Bucket
	p.Key = in.Key

}

type CompleteMultipartUploadOutput struct {

	// The name of the bucket that contains the newly created object. Does not return
	// the access point ARN or access point alias if used.
	//
	// Access points are not supported by directory buckets.
	Bucket *string

	// Indicates whether the multipart upload uses an S3 Bucket Key for server-side
	// encryption with Key Management Service (KMS) keys (SSE-KMS).
	BucketKeyEnabled *bool

	// The base64-encoded, 32-bit CRC-32 checksum of the object. This will only be
	// present if it was uploaded with the object. When you use an API operation on an
	// object that was uploaded using multipart uploads, this value may not be a direct
	// checksum value of the full object. Instead, it's a calculation based on the
	// checksum values of each individual part. For more information about how
	// checksums are calculated with multipart uploads, see [Checking object integrity]in the Amazon S3 User
	// Guide.
	//
	// [Checking object integrity]: https://docs.aws.amazon.com/AmazonS3/latest/userguide/checking-object-integrity.html#large-object-checksums
	ChecksumCRC32 *string

	// The base64-encoded, 32-bit CRC-32C checksum of the object. This will only be
	// present if it was uploaded with the object. When you use an API operation on an
	// object that was uploaded using multipart uploads, this value may not be a direct
	// checksum value of the full object. Instead, it's a calculation based on the
	// checksum values of each individual part. For more information about how
	// checksums are calculated with multipart uploads, see [Checking object integrity]in the Amazon S3 User
	// Guide.
	//
	// [Checking object integrity]: https://docs.aws.amazon.com/AmazonS3/latest/userguide/checking-object-integrity.html#large-object-checksums
	ChecksumCRC32C *string

	// The base64-encoded, 160-bit SHA-1 digest of the object. This will only be
	// present if it was uploaded with the object. When you use the API operation on an
	// object that was uploaded using multipart uploads, this value may not be a direct
	// checksum value of the full object. Instead, it's a calculation based on the
	// checksum values of each individual part. For more information about how
	// checksums are calculated with multipart uploads, see [Checking object integrity]in the Amazon S3 User
	// Guide.
	//
	// [Checking object integrity]: https://docs.aws.amazon.com/AmazonS3/latest/userguide/checking-object-integrity.html#large-object-checksums
	ChecksumSHA1 *string

	// The base64-encoded, 256-bit SHA-256 digest of the object. This will only be
	// present if it was uploaded with the object. When you use an API operation on an
	// object that was uploaded using multipart uploads, this value may not be a direct
	// checksum value of the full object. Instead, it's a calculation based on the
	// checksum values of each individual part. For more information about how
	// checksums are calculated with multipart uploads, see [Checking object integrity]in the Amazon S3 User
	// Guide.
	//
	// [Checking object integrity]: https://docs.aws.amazon.com/AmazonS3/latest/userguide/checking-object-integrity.html#large-object-checksums
	ChecksumSHA256 *string

	// Entity tag that identifies the newly created object's data. Objects with
	// different object data will have different entity tags. The entity tag is an
	// opaque string. The entity tag may or may not be an MD5 digest of the object
	// data. If the entity tag is not an MD5 digest of the object data, it will contain
	// one or more nonhexadecimal characters and/or will consist of less than 32 or
	// more than 32 hexadecimal digits. For more information about how the entity tag
	// is calculated, see [Checking object integrity]in the Amazon S3 User Guide.
	//
	// [Checking object integrity]: https://docs.aws.amazon.com/AmazonS3/latest/userguide/checking-object-integrity.html
	ETag *string

	// If the object expiration is configured, this will contain the expiration date (
	// expiry-date ) and rule ID ( rule-id ). The value of rule-id is URL-encoded.
	//
	// This functionality is not supported for directory buckets.
	Expiration *string

	// The object key of the newly created object.
	Key *string

	// The URI that identifies the newly created object.
	Location *string

	// If present, indicates that the requester was successfully charged for the
	// request.
	//
	// This functionality is not supported for directory buckets.
	RequestCharged types.RequestCharged

	// If present, indicates the ID of the KMS key that was used for object encryption.
	SSEKMSKeyId *string

	// The server-side encryption algorithm used when storing this object in Amazon S3
	// (for example, AES256 , aws:kms ).
	ServerSideEncryption types.ServerSideEncryption

	// Version ID of the newly created object, in case the bucket has versioning
	// turned on.
	//
	// This functionality is not supported for directory buckets.
	VersionId *string

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationCompleteMultipartUploadMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsRestxml_serializeOpCompleteMultipartUpload{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsRestxml_deserializeOpCompleteMultipartUpload{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, "CompleteMultipartUpload"); err != nil {
		return fmt.Errorf("add protocol finalizers: %v", err)
	}

	if err = addlegacyEndpointContextSetter(stack, options); err != nil {
		return err
	}
	if err = addSetLoggerMiddleware(stack, options); err != nil {
		return err
	}
	if err = addClientRequestID(stack); err != nil {
		return err
	}
	if err = addComputeContentLength(stack); err != nil {
		return err
	}
	if err = addResolveEndpointMiddleware(stack, options); err != nil {
		return err
	}
	if err = addComputePayloadSHA256(stack); err != nil {
		return err
	}
	if err = addRetry(stack, options); err != nil {
		return err
	}
	if err = addRawResponseToMetadata(stack); err != nil {
		return err
	}
	if err = addRecordResponseTiming(stack); err != nil {
		return err
	}
	if err = addSpanRetryLoop(stack, options); err != nil {
		return err
	}
	if err = addClientUserAgent(stack, options); err != nil {
		return err
	}
	if err = smithyhttp.AddErrorCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = addSetLegacyContextSigningOptionsMiddleware(stack); err != nil {
		return err
	}
	if err = addPutBucketContextMiddleware(stack); err != nil {
		return err
	}
	if err = addTimeOffsetBuild(stack, c); err != nil {
		return err
	}
	if err = addUserAgentRetryMode(stack, options); err != nil {
		return err
	}
	if err = addIsExpressUserAgent(stack); err != nil {
		return err
	}
	if err = addOpCompleteMultipartUploadValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opCompleteMultipartUpload(options.Region), middleware.Before); err != nil {
		return err
	}
	if err = addMetadataRetrieverMiddleware(stack); err != nil {
		return err
	}
	if err = addRecursionDetection(stack); err != nil {
		return err
	}
	if err = addCompleteMultipartUploadUpdateEndpoint(stack, options); err != nil {
		return err
	}
	if err = addResponseErrorMiddleware(stack); err != nil {
		return err
	}
	if err = v4.AddContentSHA256HeaderMiddleware(stack); err != nil {
		return err
	}
	if err = disableAcceptEncodingGzip(stack); err != nil {
		return err
	}
	if err = s3cust.HandleResponseErrorWith200Status(stack); err != nil {
		return err
	}
	if err = addRequestResponseLogging(stack, options); err != nil {
		return err
	}
	if err = addDisableHTTPSMiddleware(stack, options); err != nil {
		return err
	}
	if err = addSerializeImmutableHostnameBucketMiddleware(stack, options); err != nil {
		return err
	}
	if err = addSpanInitializeStart(stack); err != nil {
		return err
	}
	if err = addSpanInitializeEnd(stack); err != nil {
		return err
	}
	if err = addSpanBuildRequestStart(stack); err != nil {
		return err
	}
	if err = addSpanBuildRequestEnd(stack); err != nil {
		return err
	}
	return nil
}

func (v *CompleteMultipartUploadInput) bucket() (string, bool) {
	if v.Bucket == nil {
		return "", false
	}
	return *v.Bucket, true
}

func newServiceMetadataMiddleware_opCompleteMultipartUpload(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: "CompleteMultipartUpload",
	}
}

// getCompleteMultipartUploadBucketMember returns a pointer to string denoting a
// provided bucket member valueand a boolean indicating if the input has a modeled
// bucket name,
func getCompleteMultipartUploadBucketMember(input interface{}) (*string, bool) {
	in := input.(*CompleteMultipartUploadInput)
	if in.Bucket == nil {
		return nil, false
	}
	return in.Bucket, true
}
func addCompleteMultipartUploadUpdateEndpoint(stack *middleware.Stack, options Options) error {
	return s3cust.UpdateEndpoint(stack, s3cust.UpdateEndpointOptions{
		Accessor: s3cust.UpdateEndpointParameterAccessor{
			GetBucketFromInput: getCompleteMultipartUploadBucketMember,
		},
		UsePathStyle:                   options.UsePathStyle,
		UseAccelerate:                  options.UseAccelerate,
		SupportsAccelerate:             true,
		TargetS3ObjectLambda:           false,
		EndpointResolver:               options.EndpointResolver,
		EndpointResolverOptions:        options.EndpointOptions,
		UseARNRegion:                   options.UseARNRegion,
		DisableMultiRegionAccessPoints: options.DisableMultiRegionAccessPoints,
	})
}
