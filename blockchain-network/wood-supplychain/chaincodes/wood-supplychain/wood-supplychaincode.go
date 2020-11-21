package main

import (
	"encoding/json" //reading and writing JSON
	//"fmt"
	id "github.com/hyperledger/fabric-chaincode-go/pkg/cid" // import for Client Identity
	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/common/flogging"
	"strconv"
	"strings"
	"time"
)

//Logger for Logging
var logger = flogging.MustGetLogger("WoodTracker-Manager")

var OwnerNames = map[string]string{
	"forest.com":      "FOREST",
	"cutter.com":      "CUTTER",
	"logistics.com":   "TRANSPORTER",
	"manufacture.com": "MANUFACTURER",
}

type LogUpdateByCutter struct {
	TreeID        string `json:"tree_id"`
	LogDimensions string `json:"log_dimensions"`
	LogStatus     string `json:"log_status"`
	LogTime       string `json:"log_time"`
	LogLocation   string `json:"log_location"`
}

type LogUpdateByLogistics struct {
	TreeID        string `json:"tree_id"`
	LogDimensions string `json:"log_dimensions"`
	LoadingTime   string `json:"loading_time"`
	LoadingStatus string `json:"loading_status"`
}

type LogUpdateByManufacturer struct {
	TreeID            string `json:"tree_id"`
	ProductDimensions string `json:"product_dimensions"`
	Status            string `json:"status"`
	QRCode            string `json:"qr_code"`
}

type Tree struct {
	TreeID          string `json:"tree_id"`
	Owner           string `json:"owner"`
	Age             string `json:"age"`
	Quality         string `json:"quality"`
	TreeDimensions  string `json:"tree_dimensions"`
	RFIDCreatedTime string `json:"rfid_created_time"`
	Location        string `json:"location"`
	StatusOfTree    string `json:"sts_of_tree"`
}

//=========================================================
//WoodTracker structure with 6 properties , Structure tags are used by encoding/json library
//=========================================================
type WoodTrack struct {
	AssetType                   string                  `json:"asset_type"`
	TreeID                      string                  `json:"tree_id"`
	TreeDetails                 Tree                    `json:"tree"`
	TreeLogUpdateByCutter       LogUpdateByCutter       `json:"log_updated_by_cutter"`
	TreeLogUpdateByLogistics    LogUpdateByLogistics    `json:"log_updated_by_logistics"`
	TreeLogUpdateByManufacturer LogUpdateByManufacturer `json:"log_updated_by_manufacturer"`
	Creator                     string                  `json:"crtr"`
	UpdatedBy                   string                  `json:"uby"`
	CreateTs                    string                  `json:"cts"`
	UpdateTs                    string                  `json:"uts"`
}

//Smart Contract structure
type WoodTrackManager struct {
}

//Returns the complete identity in the format
//Certitificate issuer orgs's domain name
//Returns string Unkown if not able parse the invoker certificate
func (wm *WoodTrackManager) getInvokerIdentity(stub shim.ChaincodeStubInterface) (bool, string) {
	//Following id comes in the format X509::<Subject>::<Issuer>>
	enCert, err := id.GetX509Certificate(stub)
	if err != nil {
		logger.Errorf("Getting Certificate Details Failed:" + string(err.Error()))
		return false, "Unknown"
	}
	issuersOrgs := enCert.Issuer.Organization
	if len(issuersOrgs) == 0 {
		return false, "Unknown"
	}
	domainName := issuersOrgs[0]
	return true, string(domainName)
}

//Global variables used in chaincode
var jsonResp string
var errorKey string
var errorData string

//=========================================================================================================
// The Init method is called when the Smart Contract "WoodTrack" is instantiated by the blockchain network
//=========================================================================================================
func (wm *WoodTrackManager) Init(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("###### WoodTrack-Chaincode is Initialized #######")
	return shim.Success(nil)
}

//=================================================
// Invoke - Our entry point for Invocations
//==================================================
func (wm *WoodTrackManager) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	action, args := stub.GetFunctionAndParameters()
	logger.Infof("WoodTrack ChainCode is Invoked with Action Name is : " + string(action))
	switch action {
	case "addTreeDetailsByForestOfficer": // addTreeDetailsByForestOfficer
		return wm.addTreeDetailsByForestOfficer(stub, args)
	case "updateTreeDetailsByForestOfficer": // updateTreeDetailsByForestOfficer
		return wm.updateTreeDetailsByForestOfficer(stub, args)
	case "logUpdateByCutter": // logUpdateByCutter into DLT
		return wm.logUpdateByCutter(stub, args)
	case "logStatusUpdateByCutter": // logStatusUpdateByCutter into DLT
		return wm.logStatusUpdateByCutter(stub, args)
	case "loadingUpdateByLogistics": // loadingUpdateByLogistics into DLT
		return wm.loadingUpdateByLogistics(stub, args)
	case "loadingStatusUpdateByLogistics": // loadingStatusUpdateByLogistics into DLT
		return wm.loadingStatusUpdateByLogistics(stub, args)
	case "treeDetailsUpdateByManufacturer": // treeDetailsUpdateByManufacturer into DLT
		return wm.treeDetailsUpdateByManufacturer(stub, args)
	case "getTreeDetailsById": // getTreeDetailsById by ID From DLT
		return wm.getTreeDetailsById(stub, args)
	case "queryWoodTrack": //queryWoodTrack details from DLT
		return wm.queryWoodTrack(stub, args)
	default:
		logger.Errorf("Unknown Function Invoked, Available Functions : addTreeDetailsByForestOfficer,updateTreeDetailsByForestOfficer,logUpdateByCutter,loadingUpdateByLogistics,logUpdateByManufacturer,getTreeDetailsById,queryWoodTrack")
		jsonResp = "{\"Data\":" + action + ",\"ErrorDetails\":\"Available Functions:addTreeDetailsByForestOfficer,updateTreeDetailsByForestOfficer,logUpdateByCutter,loadingUpdateByLogistics,logUpdateByManufacturer,getTreeDetailsById,queryWoodTrack\"}"
		return shim.Error(jsonResp)
	}
}

// ============================================================
// addTreeDetailsByForestOfficer - create a new WoodTrack, store into chaincode state
// ============================================================
func (wm *WoodTrackManager) addTreeDetailsByForestOfficer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		logger.Errorf("addTreeDetailsByForestOfficer:Invalid Number of arguments provided for transaction")
		jsonResp = "{\"Data\":" + strconv.Itoa(len(args)) + ",\"ErrorDetails\":\"Invalid Number of argumnets provided for transaction\"}"
		return shim.Error(jsonResp)
	}
	var woodtrackObj WoodTrack
	err := json.Unmarshal([]byte(args[0]), &woodtrackObj)
	if err != nil {
		errorKey = args[0]
		errorData = "Invalid json provided as input"
		jsonResp = "{\"Data\":" + errorKey + ",\"ErrorDetails\":\"" + errorData + "\"}"
		logger.Error("addTreeDetailsByForestOfficer:" + string(jsonResp))
		return shim.Error(jsonResp)
	}
	res, certData := wm.getInvokerIdentity(stub)
	if res == false {
		return shim.Error("Unauthorized access")
	}
	if len(woodtrackObj.TreeDetails.TreeID) == 0 {
		errorKey = "undefined"
		errorData = "TreeID is Mandatory"
		jsonResp = "{\"Data\":" + errorKey + ",\"ErrorDetails\":\"" + errorData + "\"}"
		logger.Error(string(jsonResp))
		return shim.Error(jsonResp)
	}
	if certData != "forest.com" {
		errorKey = args[0]
		errorData = "Tree Details Insert by Forest Officer Only Unknown Owner is trying to Insert"
		jsonResp = "{\"Data\":" + errorKey + ",\"ErrorDetails\":\"" + errorData + "\"}"
		logger.Error("addTreeDetailsByForestOfficer:" + string(jsonResp))
		return shim.Error(jsonResp)
	}
	ownername := OwnerNames[certData]

	woodTrackExist, err := stub.GetState(woodtrackObj.TreeDetails.TreeID)
	if err != nil {
		errorKey = woodtrackObj.TreeDetails.TreeID
		replaceErr := strings.Replace(err.Error(), "\"", " ", -1)
		errorData = "GetState is Failed :" + replaceErr
		jsonResp = "{\"Data\":" + errorKey + ",\"ErrorDetails\":\"" + errorData + "\"}"
		logger.Error("addTreeDetailsByForestOfficer:" + string(jsonResp))
		return shim.Error(jsonResp)
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	if woodTrackExist == nil {
		woodtrackObj.AssetType = "wood-track"
		woodtrackObj.TreeID = woodtrackObj.TreeDetails.TreeID
		woodtrackObj.TreeDetails.Owner = ownername
		woodtrackObj.TreeDetails.StatusOfTree = "NotReady"
		woodtrackObj.Creator = certData
		woodtrackObj.UpdatedBy = certData
		woodtrackObj.CreateTs = timestamp
		woodtrackObj.UpdateTs = timestamp
		var cutter LogUpdateByCutter
		if cutter != woodtrackObj.TreeLogUpdateByCutter {

			logger.Errorf("addTreeDetailsByForestOfficer:woodtrackObj.LogUpdateByCutter it should be empty")
			jsonResp = "{\"Data\":" + woodtrackObj.TreeID + ",\"ErrorDetails\":\"woodtrackObj.LogUpdateByCutter it should be empty\"}"
			return shim.Error(jsonResp)
		}
		var logistics LogUpdateByLogistics
		if logistics != woodtrackObj.TreeLogUpdateByLogistics {
			logger.Errorf("addTreeDetailsByForestOfficer:woodtrackObj.LogUpdateByLogistics it should be empty")
			jsonResp = "{\"Data\":" + woodtrackObj.TreeID + ",\"ErrorDetails\":\"woodtrackObj.LogUpdateByLogistics it should be empty\"}"
			return shim.Error(jsonResp)
		}
		var manufacture LogUpdateByManufacturer
		if manufacture != woodtrackObj.TreeLogUpdateByManufacturer {
			logger.Errorf("addTreeDetailsByForestOfficer:woodtrackObj.LogUpdateByManufacturer it should be empty")
			jsonResp = "{\"Data\":" + woodtrackObj.TreeID + ",\"ErrorDetails\":\"woodtrackObj.LogUpdateByManufacturer it should be empty\"}"
			return shim.Error(jsonResp)
		}
		woodtrackJson, err := json.Marshal(woodtrackObj)
		if err != nil {
			logger.Errorf("addTreeDetailsByForestOfficer : Marshalling Error : " + string(err.Error()))
			replaceErr := strings.Replace(err.Error(), "\"", " ", -1)
			errorData = "Marshalling Error :" + replaceErr
			jsonResp = "{\"Data\":" + woodtrackObj.TreeID + ",\"ErrorDetails\":\"" + errorData + "\"}"
			return shim.Error(jsonResp)
		}
		err = stub.PutState(woodtrackObj.TreeID, woodtrackJson)
		if err != nil {
			logger.Errorf("addTreeDetailsByForestOfficer:PutState is Failed :" + string(err.Error()))
			jsonResp = "{\"Data\":" + woodtrackObj.TreeID + ",\"ErrorDetails\":\"Unable to set the WoodTrack\"}"
			return shim.Error(jsonResp)
		}
		logger.Infof("addTreeDetailsByForestOfficer:WoodTrack added succesfull for TreeID is :" + string(woodtrackObj.TreeID))
	} else {
		logger.Errorf("addTreeDetailsByForestOfficer:WoodTrack ID is already Exists " + string(woodtrackObj.TreeDetails.TreeID))
		jsonResp = "{\"Data\":" + woodtrackObj.TreeDetails.TreeID + ",\"ErrorDetails\":\"WoodTrack TreeID is already Exists\"}"
		return shim.Error(jsonResp)
	}

	resultData := map[string]interface{}{
		"trxnID":  stub.GetTxID(),
		"TreeID":  woodtrackObj.TreeID,
		"message": "WoodTrack added succesfull",
		"data":    woodtrackObj,
	}
	respJson, _ := json.Marshal(resultData)
	return shim.Success(respJson)
}

// ============================================================
// updateTreeDetailsByForestOfficer - update status of tree in  chaincode state
// ============================================================
func (wm *WoodTrackManager) updateTreeDetailsByForestOfficer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		logger.Errorf("updateTreeDetailsByForestOfficer:Invalid Number of arguments provided for transaction")
		jsonResp = "{\"Data\":" + strconv.Itoa(len(args)) + ",\"ErrorDetails\":\"Invalid Number of argumnets provided for transaction treeid,sts_of_tree(Ready/NotReady\"}"
		return shim.Error(jsonResp)
	}

	res, certData := wm.getInvokerIdentity(stub)
	if res == false {
		return shim.Error("Unauthorized access")
	}

	if len(args[0]) == 0 {
		errorKey = "undefined"
		errorData = "TreeId is Mandatory"
		jsonResp = "{\"Data\":" + errorKey + ",\"ErrorDetails\":\"" + errorData + "\"}"
		logger.Error(string(jsonResp))
		return shim.Error(jsonResp)
	}

	if len(args[1]) == 0 {
		errorKey = "undefined"
		errorData = "Status of Tree(Ready/NotReady) is Mandatory"
		jsonResp = "{\"Data\":" + errorKey + ",\"ErrorDetails\":\"" + errorData + "\"}"
		logger.Error(string(jsonResp))
		return shim.Error(jsonResp)
	}
	if certData != "forest.com" {
		errorKey = args[0]
		errorData = "Tree Details Update by Forest Officer Only Unknown Owner is trying to Update"
		jsonResp = "{\"Data\":" + errorKey + ",\"ErrorDetails\":\"" + errorData + "\"}"
		logger.Error("updateTreeDetailsByForestOfficer:" + string(jsonResp))
		return shim.Error(jsonResp)
	}
	ownername := OwnerNames[certData]
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	woodTrackExist, err := stub.GetState(args[0])
	if err != nil {
		errorKey = string(args[0])
		replaceErr := strings.Replace(err.Error(), "\"", " ", -1)
		errorData = "GetState is Failed :" + replaceErr
		jsonResp = "{\"Data\":" + errorKey + ",\"ErrorDetails\":\"" + errorData + "\"}"
		logger.Error("updateTreeDetailsByForestOfficer:" + string(jsonResp))
		return shim.Error(jsonResp)
	}
	if woodTrackExist == nil {
		logger.Errorf("updateTreeDetailsByForestOfficer:No Existing WoodTrack for TreeID:" + string(args[0]))
		jsonResp = "{\"Data\":\"" + args[0] + "\",\"ErrorDetails\":\"No Existing WoodTrack\"}"
		return shim.Error(jsonResp)
	} else {
		woodtrackObj := WoodTrack{}
		err := json.Unmarshal(woodTrackExist, &woodtrackObj)
		if err != nil {
			logger.Errorf("updateTreeDetailsByForestOfficer::Existing WoodTrackData unmarshalling Error" + string(err.Error()))
			replaceErr := strings.Replace(err.Error(), "\"", " ", -1)
			errorData = "Unmarshalling Error :" + replaceErr
			jsonResp = "{\"Data\":\"" + args[0] + "\",\"ErrorDetails\":\"" + errorData + "\"}"
			return shim.Error(jsonResp)
		}
		if ownername == woodtrackObj.TreeDetails.Owner {
			if args[1] == "Ready" || args[1] == "NotReady" {
				woodtrackObj.TreeDetails.StatusOfTree = args[1]
			} else {
				errorKey = args[1]
				errorData = "Invalid Status of Tree(Ready/NotReady)"
				jsonResp = "{\"Data\":" + errorKey + ",\"ErrorDetails\":\"" + errorData + "\"}"
				logger.Error(string(jsonResp))
				return shim.Error(jsonResp)
			}
			woodtrackObj.UpdateTs = timestamp

			woodtrackJson, err := json.Marshal(woodtrackObj)
			if err != nil {
				logger.Errorf("updateTreeDetailsByForestOfficer : Marshalling Error : " + string(err.Error()))
				replaceErr := strings.Replace(err.Error(), "\"", " ", -1)
				errorData = "Marshalling Error :" + replaceErr
				jsonResp = "{\"Data\":" + woodtrackObj.TreeID + ",\"ErrorDetails\":\"" + errorData + "\"}"
				return shim.Error(jsonResp)
			}
			err = stub.PutState(args[0], woodtrackJson)
			if err != nil {
				logger.Errorf("updateTreeDetailsByForestOfficer:PutState is Failed :" + string(err.Error()))
				jsonResp = "{\"Data\":" + args[0] + ",\"ErrorDetails\":\"Unable to set the WoodTrack\"}"
				return shim.Error(jsonResp)
			}
			logger.Infof("updateTreeDetailsByForestOfficer:WoodTrack added succesfull for TreeID is :" + string(args[0]))

		} else {
			logger.Errorf("updateTreeDetailsByForestOfficer:Unauthorized Organization is trying to update Existing TreeDetails")
			jsonResp = "{\"Data\":" + args[0] + ",\"ErrorDetails\":\"Access Denied for Unknown Organization\"}"
			return shim.Error(jsonResp)
		}
	}

	resultData := map[string]interface{}{
		"trxnID":  stub.GetTxID(),
		"TreeId":  args[0],
		"message": "Status Of Tree  updated Successfully.",
	}
	respJSON, _ := json.Marshal(resultData)
	return shim.Success(respJSON)

}

// ============================================================
// logUpdateByCutter - update tree details by cutter "
// ============================================================
func (wm *WoodTrackManager) logUpdateByCutter(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		logger.Errorf("logUpdateByCutter:Invalid Number of arguments provided for transaction")
		jsonResp = "{\"Data\":" + strconv.Itoa(len(args)) + ",\"ErrorDetails\":\"Invalid Number of argumnets provided for transaction\"}"
		return shim.Error(jsonResp)
	}

	res, certData := wm.getInvokerIdentity(stub)
	if res == false {
		return shim.Error("Unauthorized access")
	}

	if certData != "cutter.com" {
		errorKey = args[0]
		errorData = "Tree Details Updated by Cutter Only Unknown Owner is trying to Update"
		jsonResp = "{\"Data\":" + errorKey + ",\"ErrorDetails\":\"" + errorData + "\"}"
		logger.Error("logUpdateByCutter:" + string(jsonResp))
		return shim.Error(jsonResp)
	}
	ownername := OwnerNames[certData]

	var cutterObj LogUpdateByCutter
	err := json.Unmarshal([]byte(args[0]), &cutterObj)
	if err != nil {
		errorKey = args[0]
		errorData = "Invalid json provided as input"
		jsonResp = "{\"Data\":" + errorKey + ",\"ErrorDetails\":\"" + errorData + "\"}"
		logger.Error("logUpdateByCutter:" + string(jsonResp))
		return shim.Error(jsonResp)
	}

	if len(cutterObj.TreeID) == 0 {
		errorKey = "undefined"
		errorData = "TreeID is Mandatory"
		jsonResp = "{\"Data\":" + errorKey + ",\"ErrorDetails\":\"" + errorData + "\"}"
		logger.Error(string(jsonResp))
		return shim.Error(jsonResp)
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	woodTrackExist, err := stub.GetState(cutterObj.TreeID)
	if err != nil {
		errorKey = string(cutterObj.TreeID)
		replaceErr := strings.Replace(err.Error(), "\"", " ", -1)
		errorData = "GetState is Failed :" + replaceErr
		jsonResp = "{\"Data\":" + errorKey + ",\"ErrorDetails\":\"" + errorData + "\"}"
		logger.Error("logUpdateByCutter:" + string(jsonResp))
		return shim.Error(jsonResp)
	}
	if woodTrackExist == nil {
		logger.Errorf("logUpdateByCutter:No Existing WoodTrack for TreeID:" + string(cutterObj.TreeID))
		jsonResp = "{\"Data\":\"" + cutterObj.TreeID + "\",\"ErrorDetails\":\"No Existing WoodTrack\"}"
		return shim.Error(jsonResp)
	} else {
		woodtrackObj := WoodTrack{}
		err := json.Unmarshal(woodTrackExist, &woodtrackObj)
		if err != nil {
			logger.Errorf("logUpdateByCutter::Existing WoodTrackData unmarshalling Error" + string(err.Error()))
			replaceErr := strings.Replace(err.Error(), "\"", " ", -1)
			errorData = "Unmarshalling Error :" + replaceErr
			jsonResp = "{\"Data\":\"" + cutterObj.TreeID + "\",\"ErrorDetails\":\"" + errorData + "\"}"
			return shim.Error(jsonResp)
		}
		if woodtrackObj.TreeDetails.Owner == "FOREST" || woodtrackObj.TreeDetails.Owner=="CUTTER" {
			if woodtrackObj.TreeDetails.StatusOfTree == "Ready" {
				if cutterObj.LogStatus == "NoCut" || cutterObj.LogStatus == "Cut" {
					woodtrackObj.TreeDetails.Owner = ownername
					woodtrackObj.TreeLogUpdateByCutter.TreeID = cutterObj.TreeID
					woodtrackObj.TreeLogUpdateByCutter.LogDimensions = cutterObj.LogDimensions
					woodtrackObj.TreeLogUpdateByCutter.LogStatus = cutterObj.LogStatus
					woodtrackObj.TreeLogUpdateByCutter.LogTime = cutterObj.LogTime
					woodtrackObj.TreeLogUpdateByCutter.LogLocation = cutterObj.LogLocation
					woodtrackObj.UpdatedBy = certData
					woodtrackObj.UpdateTs = timestamp
					woodtrackJson, err := json.Marshal(woodtrackObj)
					if err != nil {
						logger.Errorf("logUpdateByCutter: Marshalling Error : " + string(err.Error()))
						replaceErr := strings.Replace(err.Error(), "\"", " ", -1)
						errorData = "Marshalling Error :" + replaceErr
						jsonResp = "{\"Data\":" + woodtrackObj.TreeID + ",\"ErrorDetails\":\"" + errorData + "\"}"
						return shim.Error(jsonResp)
					}
					err = stub.PutState(cutterObj.TreeID, woodtrackJson)
					if err != nil {
						logger.Errorf("logUpdateByCutter:PutState is Failed :" + string(err.Error()))
						jsonResp = "{\"Data\":" + cutterObj.TreeID + ",\"ErrorDetails\":\"Unable to set the WoodTrack\"}"
						return shim.Error(jsonResp)
					}
					logger.Infof("logUpdateByCutter:WoodTrack LogUpdateByCutter Updated succesfull for TreeID is :" + string(cutterObj.TreeID))
				} else {
					logger.Errorf("logUpdateByCutter:Invalid Log Status Either NoCut/Cut")
					jsonResp = "{\"Data\":" + cutterObj.LogStatus + ",\"ErrorDetails\":\"Invalid Log Status Either NoCut/Cut\"}"
					return shim.Error(jsonResp)
				}
			} else {
				logger.Errorf("logUpdateByCutter:Failed to Updated Cutter Tree is NotReady To Cut")
				jsonResp = "{\"Data\":" + cutterObj.TreeID + ",\"ErrorDetails\":\"Failed to Updated Cutter Tree is NotReady To Cut\"}"
				return shim.Error(jsonResp)
			}
		} else {
			logger.Errorf("logUpdateByCutter:Unauthorized Organization is trying to update Existing TreeDetails")
			jsonResp = "{\"Data\":" + cutterObj.TreeID + ",\"ErrorDetails\":\"Access Denied for Unknown Organization\"}"
			return shim.Error(jsonResp)

		}
	}

	resultData := map[string]interface{}{
		"trxnID":  stub.GetTxID(),
		"TreeID":  cutterObj.TreeID,
		"message": "logUpdateByCutter Updated Successfull",
		"data":    cutterObj,
	}
	respJson, _ := json.Marshal(resultData)
	return shim.Success(respJson)

}

// ============================================================
// logStatusUpdateByCutter - log status update by cutter
// ============================================================
func (wm *WoodTrackManager) logStatusUpdateByCutter(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		logger.Errorf("logStatusUpdateByCutter:Invalid Number of arguments provided for transaction")
		jsonResp = "{\"Data\":" + strconv.Itoa(len(args)) + ",\"ErrorDetails\":\"Invalid Number of argumnets provided for transaction treeid,LogStatus(NoCut/Cut)\"}"
		return shim.Error(jsonResp)
	}

	res, certData := wm.getInvokerIdentity(stub)
	if res == false {
		return shim.Error("Unauthorized access")
	}

	if len(args[0]) == 0 {
		errorKey = "undefined"
		errorData = "TreeId is Mandatory"
		jsonResp = "{\"Data\":" + errorKey + ",\"ErrorDetails\":\"" + errorData + "\"}"
		logger.Error(string(jsonResp))
		return shim.Error(jsonResp)
	}

	if len(args[1]) == 0 {
		errorKey = "undefined"
		errorData = "LogStatus (NoCut/Cut) is Mandatory"
		jsonResp = "{\"Data\":" + errorKey + ",\"ErrorDetails\":\"" + errorData + "\"}"
		logger.Error(string(jsonResp))
		return shim.Error(jsonResp)
	}
	if certData != "cutter.com" {
		errorKey = args[0]
		errorData = "Log Tree Details Update by Cutter Officer Only Unknown Owner is trying to Update"
		jsonResp = "{\"Data\":" + errorKey + ",\"ErrorDetails\":\"" + errorData + "\"}"
		logger.Error("logStatusUpdateByCutter:" + string(jsonResp))
		return shim.Error(jsonResp)
	}
	ownername := OwnerNames[certData]
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	woodTrackExist, err := stub.GetState(args[0])
	if err != nil {
		errorKey = string(args[0])
		replaceErr := strings.Replace(err.Error(), "\"", " ", -1)
		errorData = "GetState is Failed :" + replaceErr
		jsonResp = "{\"Data\":" + errorKey + ",\"ErrorDetails\":\"" + errorData + "\"}"
		logger.Error("logStatusUpdateByCutter:" + string(jsonResp))
		return shim.Error(jsonResp)
	}
	if woodTrackExist == nil {
		logger.Errorf("logStatusUpdateByCutter:No Existing WoodTrack for TreeID:" + string(args[0]))
		jsonResp = "{\"Data\":\"" + args[0] + "\",\"ErrorDetails\":\"No Existing WoodTrack\"}"
		return shim.Error(jsonResp)
	} else {
		woodtrackObj := WoodTrack{}
		err := json.Unmarshal(woodTrackExist, &woodtrackObj)
		if err != nil {
			logger.Errorf("logStatusUpdateByCutter::Existing WoodTrackData unmarshalling Error" + string(err.Error()))
			replaceErr := strings.Replace(err.Error(), "\"", " ", -1)
			errorData = "Unmarshalling Error :" + replaceErr
			jsonResp = "{\"Data\":\"" + args[0] + "\",\"ErrorDetails\":\"" + errorData + "\"}"
			return shim.Error(jsonResp)
		}
		if ownername == woodtrackObj.TreeDetails.Owner {
			if args[1] == "NoCut" || args[1] == "Cut" {
				woodtrackObj.TreeLogUpdateByCutter.LogStatus = args[1]
			} else {
				errorKey = args[1]
				errorData = "Invalid LogStatus of Tree(NoCut/Cut)"
				jsonResp = "{\"Data\":" + errorKey + ",\"ErrorDetails\":\"" + errorData + "\"}"
				logger.Error(string(jsonResp))
				return shim.Error(jsonResp)
			}
			woodtrackObj.UpdateTs = timestamp

			woodtrackJson, err := json.Marshal(woodtrackObj)
			if err != nil {
				logger.Errorf("logStatusUpdateByCutter : Marshalling Error : " + string(err.Error()))
				replaceErr := strings.Replace(err.Error(), "\"", " ", -1)
				errorData = "Marshalling Error :" + replaceErr
				jsonResp = "{\"Data\":" + woodtrackObj.TreeID + ",\"ErrorDetails\":\"" + errorData + "\"}"
				return shim.Error(jsonResp)
			}
			err = stub.PutState(args[0], woodtrackJson)
			if err != nil {
				logger.Errorf("logStatusUpdateByCutter:PutState is Failed :" + string(err.Error()))
				jsonResp = "{\"Data\":" + args[0] + ",\"ErrorDetails\":\"Unable to set the WoodTrack\"}"
				return shim.Error(jsonResp)
			}
			logger.Infof("logStatusUpdateByCutter:WoodTrack added succesfull for TreeID is :" + string(args[0]))

		} else {
			logger.Errorf("logStatusUpdateByCutter:Unauthorized Organization is trying to update Existing TreeDetails")
			jsonResp = "{\"Data\":" + args[0] + ",\"ErrorDetails\":\"Access Denied for Unknown Organization\"}"
			return shim.Error(jsonResp)
		}
	}

	resultData := map[string]interface{}{
		"trxnID":  stub.GetTxID(),
		"TreeId":  args[0],
		"message": "LogStatus Of Tree  updated Successfully.",
	}
	respJSON, _ := json.Marshal(resultData)
	return shim.Success(respJSON)

}

// ============================================================
// loadingUpdateByLogistics - update tree details by logistics "
// ============================================================
func (wm *WoodTrackManager) loadingUpdateByLogistics(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		logger.Errorf("loadingUpdateByLogistics:Invalid Number of arguments provided for transaction")
		jsonResp = "{\"Data\":" + strconv.Itoa(len(args)) + ",\"ErrorDetails\":\"Invalid Number of argumnets provided for transaction\"}"
		return shim.Error(jsonResp)
	}

	res, certData := wm.getInvokerIdentity(stub)
	if res == false {
		return shim.Error("Unauthorized access")
	}

	if certData != "logistics.com" {
		errorKey = args[0]
		errorData = "Tree Details Updated by Cutter Only Unknown Owner is trying to Update"
		jsonResp = "{\"Data\":" + errorKey + ",\"ErrorDetails\":\"" + errorData + "\"}"
		logger.Error("loadingUpdateByLogistics:" + string(jsonResp))
		return shim.Error(jsonResp)
	}
	ownername := OwnerNames[certData]

	var logisticsObj LogUpdateByLogistics
	err := json.Unmarshal([]byte(args[0]), &logisticsObj)
	if err != nil {
		errorKey = args[0]
		errorData = "Invalid json provided as input"
		jsonResp = "{\"Data\":" + errorKey + ",\"ErrorDetails\":\"" + errorData + "\"}"
		logger.Error("loadingUpdateByLogistics:" + string(jsonResp))
		return shim.Error(jsonResp)
	}

	if len(logisticsObj.TreeID) == 0 {
		errorKey = "undefined"
		errorData = "TreeID is Mandatory"
		jsonResp = "{\"Data\":" + errorKey + ",\"ErrorDetails\":\"" + errorData + "\"}"
		logger.Error(string(jsonResp))
		return shim.Error(jsonResp)
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	woodTrackExist, err := stub.GetState(logisticsObj.TreeID)
	if err != nil {
		errorKey = string(logisticsObj.TreeID)
		replaceErr := strings.Replace(err.Error(), "\"", " ", -1)
		errorData = "GetState is Failed :" + replaceErr
		jsonResp = "{\"Data\":" + errorKey + ",\"ErrorDetails\":\"" + errorData + "\"}"
		logger.Error("loadingUpdateByLogistics:" + string(jsonResp))
		return shim.Error(jsonResp)
	}
	if woodTrackExist == nil {
		logger.Errorf("loadingUpdateByLogistics:No Existing WoodTrack for TreeID:" + string(logisticsObj.TreeID))
		jsonResp = "{\"Data\":\"" + logisticsObj.TreeID + "\",\"ErrorDetails\":\"No Existing WoodTrack\"}"
		return shim.Error(jsonResp)
	} else {
		woodtrackObj := WoodTrack{}
		err := json.Unmarshal(woodTrackExist, &woodtrackObj)
		if err != nil {
			logger.Errorf("loadingUpdateByLogistics::Existing WoodTrackData unmarshalling Error" + string(err.Error()))
			replaceErr := strings.Replace(err.Error(), "\"", " ", -1)
			errorData = "Unmarshalling Error :" + replaceErr
			jsonResp = "{\"Data\":\"" + logisticsObj.TreeID + "\",\"ErrorDetails\":\"" + errorData + "\"}"
			return shim.Error(jsonResp)
		}
		if woodtrackObj.TreeDetails.Owner == "CUTTER" ||woodtrackObj.TreeDetails.Owner=="TRANSPORTER" {
			if woodtrackObj.TreeLogUpdateByCutter.LogStatus == "Cut" {
				if logisticsObj.LoadingStatus == "NotDelivered" || logisticsObj.LoadingStatus == "Delivered" {
					woodtrackObj.TreeDetails.Owner = ownername
					woodtrackObj.TreeLogUpdateByLogistics.TreeID = logisticsObj.TreeID
					woodtrackObj.TreeLogUpdateByLogistics.LogDimensions = logisticsObj.LogDimensions
					woodtrackObj.TreeLogUpdateByLogistics.LoadingTime = logisticsObj.LoadingTime
					woodtrackObj.TreeLogUpdateByLogistics.LoadingStatus = logisticsObj.LoadingStatus
					woodtrackObj.UpdatedBy = certData
					woodtrackObj.UpdateTs = timestamp
					woodtrackJson, err := json.Marshal(woodtrackObj)
					if err != nil {
						logger.Errorf("loadingUpdateByLogistics: Marshalling Error : " + string(err.Error()))
						replaceErr := strings.Replace(err.Error(), "\"", " ", -1)
						errorData = "Marshalling Error :" + replaceErr
						jsonResp = "{\"Data\":" + woodtrackObj.TreeID + ",\"ErrorDetails\":\"" + errorData + "\"}"
						return shim.Error(jsonResp)
					}
					err = stub.PutState(logisticsObj.TreeID, woodtrackJson)
					if err != nil {
						logger.Errorf("loadingUpdateByLogistics:PutState is Failed :" + string(err.Error()))
						jsonResp = "{\"Data\":" + logisticsObj.TreeID + ",\"ErrorDetails\":\"Unable to set the WoodTrack\"}"
						return shim.Error(jsonResp)
					}
					logger.Infof("loadingUpdateByLogistics:WoodTrack added succesfull for TreeID is :" + string(logisticsObj.TreeID))
				} else {
					logger.Errorf("loadingUpdateByLogistics:Invalid Loading  Status Either Delivered/NotDelivered")
					jsonResp = "{\"Data\":" + logisticsObj.TreeID + ",\"ErrorDetails\":\"Invalid Log Status Either Delivered/NotDelivered\"}"
					return shim.Error(jsonResp)
				}
			} else {
				logger.Errorf("loadingUpdateByLogistics:Failed to Updated Logistics Tree is NotReady To Delivered/NotDelivered")
				jsonResp = "{\"Data\":" + logisticsObj.TreeID + ",\"ErrorDetails\":\"Failed to Updated Cutter Tree is NotReady To Delivered/NotDelivered\"}"
				return shim.Error(jsonResp)
			}
		} else {
			logger.Errorf("loadingUpdateByLogistics:Unauthorized Organization is trying to update Existing TreeDetails")
			jsonResp = "{\"Data\":" + logisticsObj.TreeID + ",\"ErrorDetails\":\"Access Denied for Unknown Organization\"}"
			return shim.Error(jsonResp)

		}
	}

	resultData := map[string]interface{}{
		"trxnID":  stub.GetTxID(),
		"TreeID":  logisticsObj.TreeID,
		"message": "loadingUpdateByLogistics Updated Successfull",
		"data":    logisticsObj,
	}
	respJson, _ := json.Marshal(resultData)
	return shim.Success(respJson)

}

// ============================================================
// loadingStatusUpdateByLogistics - log status update by logistics
// ============================================================
func (wm *WoodTrackManager) loadingStatusUpdateByLogistics(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		logger.Errorf("loadingStatusUpdateByLogistics:Invalid Number of arguments provided for transaction")
		jsonResp = "{\"Data\":" + strconv.Itoa(len(args)) + ",\"ErrorDetails\":\"Invalid Number of argumnets provided for transaction treeid,LoadingStatus(Delivered/NotDelivered)\"}"
		return shim.Error(jsonResp)
	}

	res, certData := wm.getInvokerIdentity(stub)
	if res == false {
		return shim.Error("Unauthorized access")
	}

	if len(args[0]) == 0 {
		errorKey = "undefined"
		errorData = "TreeId is Mandatory"
		jsonResp = "{\"Data\":" + errorKey + ",\"ErrorDetails\":\"" + errorData + "\"}"
		logger.Error(string(jsonResp))
		return shim.Error(jsonResp)
	}

	if len(args[1]) == 0 {
		errorKey = "undefined"
		errorData = "LogStatus (Delivered/NotDelivered) is Mandatory"
		jsonResp = "{\"Data\":" + errorKey + ",\"ErrorDetails\":\"" + errorData + "\"}"
		logger.Error(string(jsonResp))
		return shim.Error(jsonResp)
	}
	if certData != "logistics.com" {
		errorKey = args[0]
		errorData = "Log Tree Details Update by Logistics Officer Only Unknown Owner is trying to Update"
		jsonResp = "{\"Data\":" + errorKey + ",\"ErrorDetails\":\"" + errorData + "\"}"
		logger.Error("loadingStatusUpdateByLogistics:" + string(jsonResp))
		return shim.Error(jsonResp)
	}
	ownername := OwnerNames[certData]
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	woodTrackExist, err := stub.GetState(args[0])
	if err != nil {
		errorKey = string(args[0])
		replaceErr := strings.Replace(err.Error(), "\"", " ", -1)
		errorData = "GetState is Failed :" + replaceErr
		jsonResp = "{\"Data\":" + errorKey + ",\"ErrorDetails\":\"" + errorData + "\"}"
		logger.Error("loadingStatusUpdateByLogistics:" + string(jsonResp))
		return shim.Error(jsonResp)
	}
	if woodTrackExist == nil {
		logger.Errorf("loadingStatusUpdateByLogistics:No Existing WoodTrack for TreeID:" + string(args[0]))
		jsonResp = "{\"Data\":\"" + args[0] + "\",\"ErrorDetails\":\"No Existing WoodTrack\"}"
		return shim.Error(jsonResp)
	} else {
		woodtrackObj := WoodTrack{}
		err := json.Unmarshal(woodTrackExist, &woodtrackObj)
		if err != nil {
			logger.Errorf("loadingStatusUpdateByLogistics::Existing WoodTrackData unmarshalling Error" + string(err.Error()))
			replaceErr := strings.Replace(err.Error(), "\"", " ", -1)
			errorData = "Unmarshalling Error :" + replaceErr
			jsonResp = "{\"Data\":\"" + args[0] + "\",\"ErrorDetails\":\"" + errorData + "\"}"
			return shim.Error(jsonResp)
		}
		if ownername == woodtrackObj.TreeDetails.Owner {
			if args[1] == "Delivered" || args[1] == "NotDevelivered" {
				woodtrackObj.TreeLogUpdateByLogistics.LoadingStatus = args[1]
			} else {
				errorKey = args[1]
				errorData = "Invalid LogStatus of Tree(NotDelivered/Delivered)"
				jsonResp = "{\"Data\":" + errorKey + ",\"ErrorDetails\":\"" + errorData + "\"}"
				logger.Error(string(jsonResp))
				return shim.Error(jsonResp)
			}
			woodtrackObj.UpdateTs = timestamp

			woodtrackJson, err := json.Marshal(woodtrackObj)
			if err != nil {
				logger.Errorf("loadingStatusUpdateByLogistics : Marshalling Error : " + string(err.Error()))
				replaceErr := strings.Replace(err.Error(), "\"", " ", -1)
				errorData = "Marshalling Error :" + replaceErr
				jsonResp = "{\"Data\":" + woodtrackObj.TreeID + ",\"ErrorDetails\":\"" + errorData + "\"}"
				return shim.Error(jsonResp)
			}
			err = stub.PutState(args[0], woodtrackJson)
			if err != nil {
				logger.Errorf("loadingStatusUpdateByLogistics:PutState is Failed :" + string(err.Error()))
				jsonResp = "{\"Data\":" + args[0] + ",\"ErrorDetails\":\"Unable to set the WoodTrack\"}"
				return shim.Error(jsonResp)
			}
			logger.Infof("loadingStatusUpdateByLogistics:WoodTrack loadingStatus  succesfull for TreeID is :" + string(args[0]))

		} else {
			logger.Errorf("loadingStatusUpdateByLogistics:Unauthorized Organization is trying to update Existing TreeDetails")
			jsonResp = "{\"Data\":" + args[0] + ",\"ErrorDetails\":\"Failed to Updated Cutter Tree is NotReady To Delivered/NotDelivered\"}"
			return shim.Error(jsonResp)
		}
	}

	resultData := map[string]interface{}{
		"trxnID":  stub.GetTxID(),
		"TreeId":  args[0],
		"message": "LoadingStatus Of Tree  updated Successfully.",
	}
	respJSON, _ := json.Marshal(resultData)
	return shim.Success(respJSON)

}

// ============================================================
// treeDetailsUpdateByManufacturer - update tree details by manufacture "
// ============================================================
func (wm *WoodTrackManager) treeDetailsUpdateByManufacturer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
        if len(args) != 1 {
                logger.Errorf("treeDetailsUpdateByManufacturer:Invalid Number of arguments provided for transaction")
                jsonResp = "{\"Data\":" + strconv.Itoa(len(args)) + ",\"ErrorDetails\":\"Invalid Number of argumnets provided for transaction\"}"
                return shim.Error(jsonResp)
        }

        res, certData := wm.getInvokerIdentity(stub)
        if res == false {
                return shim.Error("Unauthorized access")
        }

        if certData != "manufacture.com" {
                errorKey = args[0]
                errorData = "Tree Details Updated by Cutter Only Unknown Owner is trying to Update"
                jsonResp = "{\"Data\":" + errorKey + ",\"ErrorDetails\":\"" + errorData + "\"}"
                logger.Error("treeDetailsUpdateByManufacturer:" + string(jsonResp))
                return shim.Error(jsonResp)
        }
        ownername := OwnerNames[certData]

        var manufactureObj LogUpdateByManufacturer
        err := json.Unmarshal([]byte(args[0]), &manufactureObj)
        if err != nil {
                errorKey = args[0]
                errorData = "Invalid json provided as input"
                jsonResp = "{\"Data\":" + errorKey + ",\"ErrorDetails\":\"" + errorData + "\"}"
                logger.Error("treeDetailsUpdateByManufacturer:" + string(jsonResp))
                return shim.Error(jsonResp)
        }

        if len(manufactureObj.TreeID) == 0 {
                errorKey = "undefined"
                errorData = "TreeID is Mandatory"
                jsonResp = "{\"Data\":" + errorKey + ",\"ErrorDetails\":\"" + errorData + "\"}"
                logger.Error(string(jsonResp))
                return shim.Error(jsonResp)
        }

        timestamp := strconv.FormatInt(time.Now().Unix(), 10)
        woodTrackExist, err := stub.GetState(manufactureObj.TreeID)
        if err != nil {
                errorKey = string(manufactureObj.TreeID)
                replaceErr := strings.Replace(err.Error(), "\"", " ", -1)
                errorData = "GetState is Failed :" + replaceErr
                jsonResp = "{\"Data\":" + errorKey + ",\"ErrorDetails\":\"" + errorData + "\"}"
                logger.Error("treeDetailsUpdateByManufacturer:" + string(jsonResp))
                return shim.Error(jsonResp)
        }
        if woodTrackExist == nil {
                logger.Errorf("treeDetailsUpdateByManufacturer:No Existing WoodTrack for TreeID:" + string(manufactureObj.TreeID))
                jsonResp = "{\"Data\":\"" + manufactureObj.TreeID + "\",\"ErrorDetails\":\"No Existing WoodTrack\"}"
                return shim.Error(jsonResp)
        } else {
                woodtrackObj := WoodTrack{}
                err := json.Unmarshal(woodTrackExist, &woodtrackObj)
                if err != nil {
                        logger.Errorf("treeDetailsUpdateByManufacturer::Existing WoodTrackData unmarshalling Error" + string(err.Error()))
                        replaceErr := strings.Replace(err.Error(), "\"", " ", -1)
                        errorData = "Unmarshalling Error :" + replaceErr
                        jsonResp = "{\"Data\":\"" + manufactureObj.TreeID + "\",\"ErrorDetails\":\"" + errorData + "\"}"
                        return shim.Error(jsonResp)
                }
                if woodtrackObj.TreeDetails.Owner == "TRANSPORTER" || woodtrackObj.TreeDetails.Owner== "MANUFACTURER" {
                        if woodtrackObj.TreeLogUpdateByLogistics.LoadingStatus == "Delivered" {
                                if  manufactureObj.Status == "NotManufactured" ||  manufactureObj.Status == "Manufactured" {
                                        woodtrackObj.TreeDetails.Owner = ownername
					 woodtrackObj.TreeLogUpdateByManufacturer.TreeID =manufactureObj.TreeID
					 woodtrackObj.TreeLogUpdateByManufacturer.ProductDimensions=manufactureObj.ProductDimensions
					 woodtrackObj.TreeLogUpdateByManufacturer.Status=manufactureObj.Status
					 woodtrackObj.TreeLogUpdateByManufacturer.QRCode=manufactureObj.QRCode
                                        woodtrackObj.UpdatedBy = certData
                                        woodtrackObj.UpdateTs = timestamp
                                        woodtrackJson, err := json.Marshal(woodtrackObj)
                                        if err != nil {
                                                logger.Errorf("treeDetailsUpdateByManufacturer: Marshalling Error : " + string(err.Error()))
                                                replaceErr := strings.Replace(err.Error(), "\"", " ", -1)
                                                errorData = "Marshalling Error :" + replaceErr
                                                jsonResp = "{\"Data\":" + woodtrackObj.TreeID + ",\"ErrorDetails\":\"" + errorData + "\"}"
                                                return shim.Error(jsonResp)
                                        }
                                        err = stub.PutState(manufactureObj.TreeID, woodtrackJson)
                                        if err != nil {
                                                logger.Errorf("treeDetailsUpdateByManufacturer:PutState is Failed :" + string(err.Error()))
                                                jsonResp = "{\"Data\":" + args[0] + ",\"ErrorDetails\":\"Unable to set the WoodTrack\"}"
                                                return shim.Error(jsonResp)
                                        }
                                        logger.Infof("treeDetailsUpdateByManufacturer:WoodTrack added succesfull for TreeID is :" + string(manufactureObj.TreeID))
                                } else {
                                        logger.Errorf("treeDetailsUpdateByManufacturer:Invalid Status Either Manufactured/Manufactured")
                                        jsonResp = "{\"Data\":" + manufactureObj.TreeID + ",\"ErrorDetails\":\"Invalid Log Status Either Manufactured/NotManufactured\"}"
                                        return shim.Error(jsonResp)
                                }
                        } else {
                                logger.Errorf("treeDetailsUpdateByManufacturer:Failed to Updated Logistics Tree is NotReady To Manufactured/NotManufactured")
                                jsonResp = "{\"Data\":" + manufactureObj.TreeID + ",\"ErrorDetails\":\"Failed to Updated Cutter Tree is NotReady To Manufactured/NotManufactured\"}"
                                return shim.Error(jsonResp)
                        }
                } else {
                        logger.Errorf("treeDetailsUpdateByManufacturer:Unauthorized Organization is trying to update Existing TreeDetails")
                        jsonResp = "{\"Data\":" + manufactureObj.TreeID + ",\"ErrorDetails\":\"Only Manufacture can update\"}"
                        return shim.Error(jsonResp)

                }
        }

        resultData := map[string]interface{}{
                "trxnID":  stub.GetTxID(),
                "TreeID":  manufactureObj.TreeID,
                "message": "treeDetailsUpdateByManufacturer Updated Successfull",
                "data":    manufactureObj,
        }
        respJson, _ := json.Marshal(resultData)
        return shim.Success(respJson)

}




// ===============================================
// getTreeDetailsById -  read a tree details from chaincode state
// ===============================================
func (wm *WoodTrackManager) getTreeDetailsById(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		logger.Errorf("getTreeDetailsById:Invalid Number of arguments are provided for transaction")
		jsonResp = "{\"Data\":" + strconv.Itoa(len(args)) + ",\"ErrorDetails\":\"Invalid Number of argumnets provided for transaction\"}"
		return shim.Error(jsonResp)
	}
	if len(args[0]) == 0 {
		errorKey = "undefined"
		errorData = "TreeId is Mandatory"
		jsonResp = "{\"Data\":" + errorKey + ",\"ErrorDetails\":\"" + errorData + "\"}"
		logger.Error(string(jsonResp))
		return shim.Error(jsonResp)
	}

	var records []WoodTrack
	woodTrackExist, err := stub.GetState(args[0])
	if err != nil {
		errorKey = args[0]
		replaceErr := strings.Replace(err.Error(), "\"", " ", -1)
		errorData = "GetState is Failed :" + replaceErr
		jsonResp = "{\"Data\":\"" + errorKey + "\",\"ErrorDetails\":\"" + errorData + "\"}"
		logger.Error("getTreeDetailsById:" + string(jsonResp))
		return shim.Error(string(jsonResp))
	}
	if woodTrackExist == nil {
		logger.Errorf("getTreeDetailsById:No Existing WoodTrack for TreeID:" + string(args[0]))
		jsonResp = "{\"Data\":\"" + args[0] + "\",\"ErrorDetails\":\"No Existing WoodTrack\"}"
		return shim.Error(jsonResp)
	} else {
		woodtrack := WoodTrack{}
		err := json.Unmarshal(woodTrackExist, &woodtrack)
		if err != nil {
			logger.Errorf("getTreeDetailsById::Existing WoodTrackData unmarshalling Error" + string(err.Error()))
			replaceErr := strings.Replace(err.Error(), "\"", " ", -1)
			errorData = "Unmarshalling Error :" + replaceErr
			jsonResp = "{\"Data\":\"" + args[0] + "\",\"ErrorDetails\":\"" + errorData + "\"}"
			return shim.Error(jsonResp)
		}
		records = append(records, woodtrack)

	}
	resultData := map[string]interface{}{
		"status": "true",
		"data":   records[0],
	}

	respJson, _ := json.Marshal(resultData)
	return shim.Success(respJson)
}

func (wm *WoodTrackManager) queryWoodTrack(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		logger.Errorf("queryWoodTrack:Invalid number of arguments are provided for transaction")
		jsonResp = "{\"Data\":" + strconv.Itoa(len(args)) + ",\"ErrorDetails\":\"Invalid Number of argumnets provided for transaction\"}"
		return shim.Error(jsonResp)
	}
	if len(args[0]) == 0 {
		errorKey = "undefined"
		errorData = "Selector Query is Mandatory"
		jsonResp = "{\"Data\":" + errorKey + ",\"ErrorDetails\":\"" + errorData + "\"}"
		logger.Error(string(jsonResp))
		return shim.Error(jsonResp)
	}

	var records []WoodTrack
	queryString := args[0]
	logger.Infof("Query Selector : " + string(queryString))
	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		logger.Error("queryWoodTrack:GetQueryResult is Failed with error :" + string(err.Error()))
		errorData = "GetQueryResult Error :" + string(err.Error())
		jsonResp = "{\"Data\":" + args[0] + ",\"ErrorDetails\":\"" + errorData + "\"}"
		return shim.Error(jsonResp)
	}
	for resultsIterator.HasNext() {
		record := WoodTrack{}
		recordBytes, _ := resultsIterator.Next()
		if (string(recordBytes.Value)) == "" {
			continue
		}
		err = json.Unmarshal(recordBytes.Value, &record)
		if err != nil {
			logger.Errorf("queryWoodTrack:Unable to unmarshal WoodTrack retrieved :" + string(err.Error()))
			replaceErr := strings.Replace(err.Error(), "\"", " ", -1)
			errorData = "Unmarshalling Error :" + replaceErr
			jsonResp = "{\"Data\":" + string(recordBytes.Value) + ",\"ErrorDetails\":\"" + errorData + "\"}"
			return shim.Error(jsonResp)
		}
		records = append(records, record)
	}
	var status = map[string]string{
		"code":            "200",
		"error":           "null",
		"display_message": "success",
		"debug_message":   "null",
	}

	resultData := map[string]interface{}{
		"status": status,
		"data":   records,
	}

	respJson, _ := json.Marshal(resultData)
	return shim.Success(respJson)
}

// ===================================================================================
//main function for the WoodTrack ChainCode
// ===================================================================================
func main() {
	err := shim.Start(new(WoodTrackManager))
	if err != nil {
		logger.Error("Error Starting WoodTrackManager Chaincode is " + string(err.Error()))
	} else {
		logger.Info("Starting WoodTrackManager Chaincode")
	}
}
