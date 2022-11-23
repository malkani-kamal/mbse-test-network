package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type Details struct {
	ModelId   string `json:"modelId"`
	Version   string `json:"version"`
	ModelName string `json:"modelName"`
	ProjectId string `json:"projectId"`
	ModelType string `json:"modelType"`
	AssetType string `json:"assetType"`
}

type Summary struct {
	SummaryId string `json:"summaryId"`
	Version   string `json:"version"`
	AssetType string `json:"assetType"`
	ModelId   string `json:"modelId"`
	ProjectId string `json:"projectId"`
	ModelName string `json:"modelName"`
	ModelType string `json:"modelType"`
	Reviewer  string `json:"reviewer"`
	Approver  string `json:"approver"`
	Comments  string `json:"comments"`
}

// GetState(id)
//PutState(id, data)
//PutState(id, data )

func (contract *SmartContract) CreateDetailedMBSEModelPrivate(ctx contractapi.TransactionContextInterface, MBSEData string) (err error) {

	fmt.Printf("CreateDetailedMBSEModelPrivate start-->")

	userRole, err := contract.GetUserRoles(ctx)
	if err != nil {
		return fmt.Errorf("User role not defined properly. %s", err.Error())
	}

	orgRole, er1 := contract.GetOrgRoles(ctx)
	if er1 != nil {
		return fmt.Errorf("Org role not defined properly. %s", er1.Error())
	}

	if orgRole != "DEVELOPER" && userRole != "CSE" {
		return fmt.Errorf("Insufficient Roles! DEVELOPER & CSE roles are required.")
	}

	stub := ctx.GetStub()

	var detail Details

	err1 := json.Unmarshal([]byte(MBSEData), &detail)
	if err1 != nil {
		return fmt.Errorf("Failed to parse MBSE argument. %s", err1.Error())
	}

	detailAsBytes, err := ctx.GetStub().GetPrivateData("_implicit_org_Org1MSP", detail.ModelId)
	if err != nil {
		return fmt.Errorf("Failed to get detail report:" + err.Error())
	} else if detailAsBytes != nil {
		return fmt.Errorf("Detail MBSE alread exist: " + detail.ModelId)
	}

	err = stub.PutPrivateData("_implicit_org_Org1MSP", detail.ModelId, []byte(MBSEData))
	if err != nil {
		return fmt.Errorf("Failed to insert MBSE in ledger. %s", err.Error())
	}
	fmt.Println("MBSE prviate model created successfully")
	return nil
}

func (contract *SmartContract) UpdateDetailedMBSEModelPrivate(ctx contractapi.TransactionContextInterface, MBSEData string) (err error) {

	fmt.Printf("UpdateDetailedMBSEModelPrivate new start-->")

	userRole, err := contract.GetUserRoles(ctx)
	if err != nil {
		return fmt.Errorf("User role not defined properly. %s", err.Error())
	}

	orgRole, er1 := contract.GetOrgRoles(ctx)
	if er1 != nil {
		return fmt.Errorf("Org role not defined properly. %s", er1.Error())
	}

	if orgRole != "DEVELOPER" && userRole != "CSE" {
		return fmt.Errorf("Insufficient Roles! DEVELOPER & CSE roles are required.")
	}

	stub := ctx.GetStub()

	var detail Details

	err1 := json.Unmarshal([]byte(MBSEData), &detail)
	if err1 != nil {
		return fmt.Errorf("Failed to parse MBSE argument. %s", err1.Error())
	}

	detailAsBytes, err := ctx.GetStub().GetPrivateData("_implicit_org_Org1MSP", detail.ModelId)
	if err != nil {
		return fmt.Errorf("Failed to get detail report:" + err.Error())
	} else if detailAsBytes == nil {
		return fmt.Errorf("Detail MBSE does not exist: " + detail.ModelId)
	}

	detailToUpdate := Details{}
	err = json.Unmarshal(detailAsBytes, &detailToUpdate) //unmarshal it aka JSON.parse()
	if err != nil {
		return fmt.Errorf("failed to unmarshal Summary JSON: %s", err.Error())
	}

	detailToUpdate.ModelName = detail.ModelName //change the owner
	detailToUpdate.ProjectId = detail.ProjectId //change the owner

	detailJSONasBytes, _ := json.Marshal(detailToUpdate)

	err = stub.PutPrivateData("_implicit_org_Org1MSP", detail.ModelId, detailJSONasBytes)
	if err != nil {
		return fmt.Errorf("Failed to insert MBGetStateSE in ledger. %s", err.Error())
	}
	fmt.Println("MBSE prviate model created successfully")
	return nil
}

func (contract *SmartContract) DeleteDetailedMBSEModelPrivate(ctx contractapi.TransactionContextInterface, mbse_id string) (err error) {

	fmt.Printf("DeleteDetailedMBSEModelPrivate start-->", mbse_id)

	// userRole, err := contract.GetUserRoles(ctx)
	// if err != nil {
	// 	return fmt.Errorf("User role not defined properly. %s", err.Error())
	// }

	// orgRole, er1 := contract.GetOrgRoles(ctx)
	// if er1 != nil {
	// }
	// 	return fmt.Errorf("Org role not defined properly. %s", er1.Error())

	// if orgRole != "DEVELOPER" && userRole != "CSE" {
	// 	return fmt.Errorf("Insufficient Roles! DEVELOPER & CSE roles are required.")
	// }

	// stub := ctx.GetStub()

	mbseBytes, err1 := ctx.GetStub().GetPrivateData("_implicit_org_Org1MSP", mbse_id)
	if err1 != nil {
		return fmt.Errorf("Failed to get the MBSE. %s", err1.Error())
	}

	if mbseBytes == nil {
		return fmt.Errorf("No MBSE record found with mbse_id. %s", mbse_id)
	}

	err = ctx.GetStub().DelPrivateData("_implicit_org_Org1MSP", mbse_id)
	if err != nil {
		return fmt.Errorf("Failed to delete MBSE in ledger. %s", err.Error())
	}

	fmt.Println("MBSE prviate model deleted successfully")
	return nil

}

//Function to get extract the userId from ca identity.  It is required to for checking the minter
func (contract *SmartContract) GetOrgRoles(ctx contractapi.TransactionContextInterface) (string, error) {

	fmt.Printf("GetOrgRoles start-->")

	orgRole, found, _ := ctx.GetClientIdentity().GetAttributeValue("orgRole")
	if found == false {
		fmt.Println("orgRole not found!")
		return "", fmt.Errorf("orgRole not found!")
	}

	return orgRole, nil
}

//Function to get extract the userId from ca identity.  It is required to for checking the minter
func (contract *SmartContract) GetOrgName(ctx contractapi.TransactionContextInterface) (string, error) {

	fmt.Println("GetOrgName start-->")

	orgName, found, _ := ctx.GetClientIdentity().GetAttributeValue("organizationName")
	if found == false {
		fmt.Println("orgName not found!")
		return "", fmt.Errorf("orgName not found!")
	}

	return orgName, nil
}

//Function to get extract the userId from ca identity.  It is required to for checking the minter
func (contract *SmartContract) GetUserRoles(ctx contractapi.TransactionContextInterface) (string, error) {

	fmt.Printf("GetUserRoles start-->")

	userRole, ufound, _ := ctx.GetClientIdentity().GetAttributeValue("userRole")
	if ufound == false {
		return "", fmt.Errorf("userRole not found!")
	}

	return userRole, nil
}

//new
//Function to get extract the userId from ca identity.  It is required to for checking the minter
func (contract *SmartContract) GetMyMSPID(ctx contractapi.TransactionContextInterface) (string, error) {

	fmt.Printf("GetMyMSPID start-->")

	uMSPId, ufound := ctx.GetClientIdentity().GetMSPID()
	if ufound == nil {
		return "", fmt.Errorf("MSP ID not found!")
	}

	fmt.Println("uMSPId:", uMSPId)
	return uMSPId, nil
}

func (contract *SmartContract) CreateSummaryMBSEModel(ctx contractapi.TransactionContextInterface, MBSEData string) (err error) {

	fmt.Printf("CreateDetailedMBSEModel start-->")

	stub := ctx.GetStub()
	var summary Summary

	err1 := json.Unmarshal([]byte(MBSEData), &summary)
	if err1 != nil {
		return fmt.Errorf("Failed to parse MBSE argument. %s", err1.Error())
	}

	err = stub.PutState(summary.SummaryId, []byte(MBSEData))
	if err != nil {
		return fmt.Errorf("Failed to insert MBSE in ledger. %s", err.Error())
	}
	fmt.Println("MBSE model created successfully")

	return nil
}

//new
func (contract *SmartContract) UpdateSummaryMBSEModel(ctx contractapi.TransactionContextInterface, MBSEData string) (err error) {

	fmt.Printf("UpdateSummaryMBSEModel start-->")

	stub := ctx.GetStub()
	var summary Summary

	err1 := json.Unmarshal([]byte(MBSEData), &summary)
	if err1 != nil {
		return fmt.Errorf("Failed to parse MBSE argument. %s", err1.Error())
	}

	summaryAsBytes, err := ctx.GetStub().GetState(summary.SummaryId)

	if err != nil {
		return fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if summaryAsBytes == nil {
		return fmt.Errorf("Summary %s does not exist", summary.SummaryId)
	}

	summaryData := new(Summary)
	_ = json.Unmarshal(summaryAsBytes, summaryData)

	summaryData.ProjectId = summary.ProjectId
	summaryData.Reviewer = summary.Reviewer
	summaryData.Approver = summary.Approver
	summaryData.Comments = summary.Comments

	summaryDAsBytes, _ := json.Marshal(summaryData)

	err = stub.PutState(summary.SummaryId, summaryDAsBytes)
	if err != nil {
		return fmt.Errorf("Failed to update MBSE summary in ledger. %s", err.Error())
	}
	fmt.Println("MBSE model updated successfully")

	return nil
}

//new
func (contract *SmartContract) GetDetailedMBSEModelPrivate(ctx contractapi.TransactionContextInterface, mbse_id string) (*Details, error) {

	fmt.Printf("GetDetailedMBSEModelPrivate start-->", mbse_id)

	myMSPID, _ := contract.GetOrgName(ctx)

	fmt.Println("GetDetailedMBSEModelPrivate myMSPID-->", myMSPID)
	// userRole, err := contract.GetUserRoles(ctx)
	// if err != nil {
	// 	return fmt.Errorf("User role not defined properly. %s", err.Error())
	// }

	// orgRole, er1 := contract.GetOrgRoles(ctx)
	// if er1 != nil {
	// }
	// 	return fmt.Errorf("Org role not defined properly. %s", er1.Error())

	// if orgRole != "DEVELOPER" && userRole != "CSE" {
	// 	return fmt.Errorf("Insufficient Roles! DEVELOPER & CSE roles are required.")
	// }

	// stub := ctx.GetStub()

	mbseBytes, err1 := ctx.GetStub().GetPrivateData("_implicit_org_Org1MSP", mbse_id)
	if err1 != nil {
		return nil, fmt.Errorf("Failed to get the MBSE. %s", err1.Error())
	}

	if mbseBytes == nil {
		return nil, fmt.Errorf("No MBSE record found with mbse_id. %s", mbse_id)
	}

	// var details Details
	// errs := json.Unmarshal(mbseBytes, details)
	// if errs != nil {
	// 	fmt.Errorf("Error in parsing the account data - %s", errs.Error())
	// 	return nil, errs
	// }

	fmt.Println("MBSE prviate model fetched successfully")

	details := new(Details)
	_ = json.Unmarshal(mbseBytes, details)

	return details, nil

}

// ===========================================================
// transfer a report by setting a new owner name on the marble
// ===========================================================
func (contract *SmartContract) ShareReport(ctx contractapi.TransactionContextInterface, mbse_id string) error {

	fmt.Println("ShareReport ->>>")

	mbseBytes, err1 := ctx.GetStub().GetPrivateData("_implicit_org_Org1MSP", mbse_id)
	if err1 != nil {
		return fmt.Errorf("Failed to get the MBSE. %s", err1.Error())
	}

	if mbseBytes == nil {
		return fmt.Errorf("No MBSE record found with mbse_id. %s", mbse_id)
	}

	fmt.Println("MBSE prviate model fetched successfully")

	details := new(Details)
	_ = json.Unmarshal(mbseBytes, details)

	fmt.Println("details", details)

	var summary Summary

	summary.ModelId = details.ModelId
	summary.ModelName = details.ModelName
	summary.ModelType = details.ModelType
	summary.AssetType = "Summary"
	summary.SummaryId = "Summary_" + details.ModelId

	fmt.Println("summaryBytes ->>>", summary)

	summaryBytes, _ := json.Marshal(summary)
	fmt.Println("summaryBytes after marshal->>>", summaryBytes)

	err := ctx.GetStub().PutState(summary.SummaryId, summaryBytes) //rewrite the marble
	if err != nil {
		return err
	}

	return nil

}

// func (contract *SmartContract) UpdateDetailedMBSEModelPrivate(ctx contractapi.TransactionContextInterface, MBSEData string) (err error) {

// 	fmt.Printf("UpdateDetailedMBSEModelPrivate new start-->")

// 	userRole, err := contract.GetUserRoles(ctx)
// 	if err != nil {
// 		return fmt.Errorf("User role not defined properly. %s", err.Error())
// 	}

// 	orgRole, er1 := contract.GetOrgRoles(ctx)
// 	if er1 != nil {
// 		return fmt.Errorf("Org role not defined properly. %s", er1.Error())
// 	}

// 	if orgRole != "DEVELOPER" && userRole != "CSE" {
// 		return fmt.Errorf("Insufficient Roles! DEVELOPER & CSE roles are required.")
// 	}

// 	stub := ctx.GetStub()

// 	var summary Summary

// 	err1 := json.Unmarshal([]byte(MBSEData), &summary)
// 	if err1 != nil {
// 		return fmt.Errorf("Failed to parse MBSE argument. %s", err1.Error())
// 	}

// 	summaryAsBytes, err := ctx.GetStub().GetPrivateData("_implicit_org_Org1MSP", summary.ModelId)
// 	if err != nil {
// 		return fmt.Errorf("Failed to get Summary report:" + err.Error())
// 	} else if summaryAsBytes == nil {
// 		return fmt.Errorf("Marble does not exist: " + summary.ModelId)
// 	}

// 	summaryToUpdate := Summary{}
// 	err = json.Unmarshal(summaryAsBytes, &summaryToUpdate) //unmarshal it aka JSON.parse()
// 	if err != nil {
// 		return fmt.Errorf("failed to unmarshal Summary JSON: %s", err.Error())
// 	}

// 	summaryToUpdate.ModelName = summary.ModelName //change the owner
// 	summaryToUpdate.ProjectId = summary.ProjectId //change the owner

// 	summaryJSONasBytes, _ := json.Marshal(summaryToUpdate)

// 	err = stub.PutPrivateData("_implicit_org_Org1MSP", summary.ModelId, summaryJSONasBytes)
// 	if err != nil {
// 		return fmt.Errorf("Failed to insert MBGetStateSE in ledger. %s", err.Error())
// 	}
// 	fmt.Println("MBSE prviate model created successfully")
// 	return nil
// }

// func (contract *SmartContract) CreateSummaryMBSEModel(ctx contractapi.TransactionContextInterface, MBSEData string) (err error) {

// 	fmt.Printf("CreateDetailedMBSEModel start-->")

// 	stub := ctx.GetStub()
// 	var detail Details

// 	err1 := json.Unmarshal([]byte(MBSEData), &detail)
// 	if err1 != nil {
// 		return fmt.Errorf("Failed to parse MBSE argument. %s", err1.Error())
// 	}

// 	err = stub.PutState(detail.ModelId, []byte(MBSEData))
// 	if err != nil {
// 		return fmt.Errorf("Failed to insert MBSE in ledger. %s", err.Error())
// 	}
// 	fmt.Println("MBSE model created successfully")

// 	return nil
// }

// //Function to approve the hashlock
// func (contract *SmartContract) Approve(ctx contractapi.TransactionContextInterface, lock_id string) (string, error) {

// 	hashLockBytes, err := ctx.GetStub().GetState(lock_id)

// 	if err != nil {
// 		fmt.Errorf("Error in fecthing hashLock from world state - %s", err.Error())
// 		return "fail", fmt.Errorf("Error in fecthing hashLock from world state - %s", err.Error())
// 	}
// 	// fmt.Printf("accountBytes: %s", accountBytes)

// 	if hashLockBytes == nil {
// 		fmt.Errorf("HashLock does not exist with ID - %s", lock_id)
// 		return "fail", fmt.Errorf("HashLock does not exist with ID - %s", lock_id)

// 	}

// 	var hashlock HashTimeLock
// 	errs := json.Unmarshal(hashLockBytes, &hashlock)
// 	if errs != nil {
// 		fmt.Errorf("Error in parsing the hashlock data - %s", errs.Error())
// 		return "fail", errs
// 	}

// 	// Approving Hashlock
// 	hashlock.Approved = 0

// 	hashlockTxn, err := json.Marshal(hashlock)

// 	err = ctx.GetStub().PutState(lock_id, hashlockTxn)
// 	if err != nil {
// 		return "fail", fmt.Errorf("Failed to insert token in ledger. %s", err.Error())
// 	}

// 	return "Hash lock approved successfully!", nil
// }

// //Function to get extract the userId from ca identity.  It is required to for checking the minter
// func (contract *SmartContract) getUserId(ctx contractapi.TransactionContextInterface) (string, error) {

// 	fmt.Printf("getUserId start-->")

// 	b64ID, err := ctx.GetClientIdentity().GetID()
// 	if err != nil {
// 		return "", fmt.Errorf("Failed to read clientID: %v", err)
// 	}
// 	decodeID, err := base64.StdEncoding.DecodeString(b64ID)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to base64 decode clientID: %v", err)
// 	}

// 	fmt.Println("minter: %s", string(decodeID))

// 	completeId := string(decodeID)
// 	userId := completeId[(strings.Index(completeId, "x509::CN=") + 9):strings.Index(completeId, ",")]
// 	fmt.Println("userId:----------", userId)

// 	return userId, nil
// }

// //MintToken will add the token to Minter's account if token's are available in Token Definition.
// //It will also reduce the number of tokens from Token Definition.
// func (contract *SmartContract) MintToken(ctx contractapi.TransactionContextInterface, transaction string) (string, error) {
// 	fmt.Printf("MintToken: %s", transaction)

// 	var transactionData Transaction
// 	errs := json.Unmarshal([]byte(transaction), &transactionData)
// 	if errs != nil {
// 		return "", fmt.Errorf("Failed to parse transaction data.", errs.Error())
// 	}

// 	minter, _ := contract.getUserId(ctx)

// 	var account Accounts
// 	account.AccountId = transactionData.ToAccountId
// 	account.TokenId = transactionData.TokenId
// 	account.Minter = minter
// 	key, _ := ctx.GetStub().CreateCompositeKey("account", []string{account.AccountId, account.TokenId})

// 	//checking the existing token balance and increasing if tokens are already available
// 	balance, _ := contract.GetBalance(ctx, key)
// 	fmt.Println("balance-", balance)
// 	if balance == -1 {
// 		account.BalanceTokens = transactionData.Amount
// 	} else {
// 		account.BalanceTokens = transactionData.Amount + balance
// 	}

// 	fmt.Println(transactionData.FromAccountId, balance)

// 	//Checking if Token definition have token supply available
// 	var token Tokens
// 	tokenAccount, err := ctx.GetStub().GetState(account.TokenId)
// 	err = json.Unmarshal(tokenAccount, &token)
// 	if err != nil {
// 		fmt.Errorf("Error in parsing the token data - %s", err.Error())
// 		return "fail", err
// 	}
// 	if (token.TokenSupply - token.TokenIssued) < transactionData.Amount {
// 		fmt.Errorf("Insufficient token!")
// 		return "fail", fmt.Errorf("Insufficient token!")
// 	} else {
// 		token.TokenIssued += transactionData.Amount
// 	}

// 	//Adding token to Minters account balance
// 	accountTxn, err := json.Marshal(account)
// 	err = ctx.GetStub().PutState(key, accountTxn)
// 	if err != nil {
// 		return "", fmt.Errorf("Error while adding transaction to ledger - %s", err.Error())
// 	}

// 	//Reducing token from token supply.
// 	tokenTxn, err := json.Marshal(token)
// 	err = ctx.GetStub().PutState(account.TokenId, tokenTxn)
// 	if err != nil {
// 		return "", fmt.Errorf("Error while adding transaction to ledger - %s", err.Error())
// 	}

// 	return "Tokens reduced from token balance.", nil
// }

// //Function to check the Toke balance of particular account.
// func (contract *SmartContract) GetBalance(ctx contractapi.TransactionContextInterface, id string) (int, error) {

// 	fmt.Printf("GetBalance: %s", id)

// 	accountBytes, err := ctx.GetStub().GetState(id)

// 	if err != nil {
// 		fmt.Errorf("Error in fecthing account balance from world state - %s", err.Error())
// 		return -1, fmt.Errorf("Error in fetching account from world state - %s", err.Error())
// 	}
// 	// fmt.Printf("accountBytes: %s", accountBytes)

// 	if accountBytes == nil {
// 		fmt.Errorf("Account does not exist with ID - %s", id)
// 		return -1, nil
// 	}

// 	var account Accounts
// 	errs := json.Unmarshal(accountBytes, &account)
// 	if errs != nil {
// 		fmt.Errorf("Error in parsing the account data - %s", errs.Error())
// 		return -1, errs
// 	}

// 	return account.BalanceTokens, nil
// }

// //Function to check if current user is minter?
// func (contract *SmartContract) isMinter(ctx contractapi.TransactionContextInterface, id string, currUser string) (int, error) {

// 	fmt.Printf("isMinter: %s", id)

// 	accountBytes, err := ctx.GetStub().GetState(id)
// 	if err != nil {
// 		fmt.Errorf("Error in fecthing account from world state - %s", err.Error())
// 		return -1, fmt.Errorf("Error in fecthing account from world state - %s", err.Error())
// 	}
// 	// fmt.Printf("accountBytes: %s", accountBytes)

// 	if accountBytes == nil {
// 		fmt.Errorf("Account does not exist with ID - %s", id)
// 		return -1, nil
// 	}

// 	var account Accounts
// 	errs := json.Unmarshal(accountBytes, &account)
// 	if errs != nil {
// 		fmt.Errorf("Error in parsing the account data - %s", errs.Error())
// 		return -1, errs
// 	}

// 	if account.Minter == currUser {
// 		return 0, nil
// 	} else {
// 		return 1, nil
// 	}

// }

// //Tansferring token from one account to other account
// func (contract *SmartContract) transactToken(ctx contractapi.TransactionContextInterface, key string, amount int, action string, toAccount string, tokenId string) (string, error) {

// 	fmt.Printf("transactToken-->", key, amount)

// 	var account Accounts
// 	accountData, _ := ctx.GetStub().GetState(key)
// 	err := json.Unmarshal(accountData, &account)
// 	if err != nil {
// 		fmt.Errorf("Error in parsing the account data - %s", err.Error())
// 		return "fail", err
// 	}

// 	//If operator is + then adding the balance else reducing the balance.
// 	if action == "+" {
// 		account.BalanceTokens += amount
// 	} else {
// 		fmt.Println(" -  operator.")

// 		//Checking the balance before reducing the number of tokens
// 		if account.BalanceTokens < amount {
// 			fmt.Errorf("Insufficient balance in account!")
// 			return "fail", fmt.Errorf("Insufficient balance in account!")
// 		}
// 		account.BalanceTokens -= amount
// 	}

// 	//Updating the token balance in the ledger.
// 	accountTxn, _ := json.Marshal(account)
// 	err = ctx.GetStub().PutState(key, accountTxn)
// 	if err != nil {
// 		return "", fmt.Errorf("Error while burn token to ledger - %s", err.Error())
// 	}

// 	fmt.Println("Balance updated successfully.")
// 	return "success", nil
// }

// func (contract *SmartContract) CreateAccount(ctx contractapi.TransactionContextInterface, accountData string) (string, error) {
// 	fmt.Printf("createAccount start-->")

// 	stub := ctx.GetStub()
// 	var account Accounts

// 	err1 := json.Unmarshal([]byte(accountData), &account)
// 	if err1 != nil {
// 		return "fail", fmt.Errorf("Failed to parse account argument. %s", err1.Error())
// 	}

// 	Key, _ := ctx.GetStub().CreateCompositeKey("account", []string{account.AccountId, account.TokenId})

// 	err := stub.PutState(Key, []byte(accountData))
// 	if err != nil {
// 		return "fail", fmt.Errorf("Failed to insert account in ledger. %s", err.Error())
// 	}
// 	fmt.Println("Account created successfully")

// 	return "success", nil
// }

// //Simple transfer function.
// func (contract *SmartContract) Transfer(ctx contractapi.TransactionContextInterface, transaction string) (string, error) {
// 	fmt.Println("Transfer----")

// 	var transactionData Transaction

// 	errs := json.Unmarshal([]byte(transaction), &transactionData)
// 	if errs != nil {
// 		return "fail", nil
// 	}

// 	fromKey, _ := ctx.GetStub().CreateCompositeKey("account", []string{transactionData.FromAccountId, transactionData.TokenId})
// 	toKey, _ := ctx.GetStub().CreateCompositeKey("account", []string{transactionData.ToAccountId, transactionData.TokenId})

// 	// Validate accounts
// 	accountBytes, err := ctx.GetStub().GetState(fromKey)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to read FromAccount from world state: %v", err)
// 	}
// 	if accountBytes == nil {
// 		return "", fmt.Errorf("FromAccount %s does not exist, Please create account first.", transactionData.FromAccountId)
// 	}

// 	toAccountBytes, err1 := ctx.GetStub().GetState(toKey)
// 	if err1 != nil {
// 		return "", fmt.Errorf("failed to read ToAccount from world state: %v", err1)
// 	}
// 	if toAccountBytes == nil {
// 		return "", fmt.Errorf("ToAccount %s does not exist, Please create account first.", transactionData.ToAccountId)
// 	}

// 	fmt.Println("Transfer----validating from key")

// 	//creating composite key AccountId+TokenId
// 	//First reducting the balance from From Account
// 	msg, err := contract.transactToken(ctx, fromKey, transactionData.Amount, "-", "", "")
// 	if err != nil {
// 		//Return incase insufficient balance error or any other error while reducing the balance
// 		return msg, fmt.Errorf("Error while reducing balance %s", err.Error())
// 	}

// 	fmt.Println("Transfer----validatedfrom key")

// 	//Adding the balance to To Account
// 	msg, err = contract.transactToken(ctx, toKey, transactionData.Amount, "+", transactionData.ToAccountId, transactionData.TokenId)
// 	if err != nil {
// 		//Return incase insufficient balance error or any other error while reducing the balance
// 		return msg, fmt.Errorf("Error while addbalance %s", err.Error())
// 	}
// 	return "success", nil
// }

// //Burn token from accounts
// func (contract *SmartContract) Burn(ctx contractapi.TransactionContextInterface, FromAccountId string, TokenId string, Amount string) (string, error) {
// 	fmt.Printf("Burn-->", FromAccountId, TokenId, Amount)
// 	burnAmount, _ := strconv.Atoi(Amount)
// 	//creating composite key AccountId+TokenId
// 	fromKey, _ := ctx.GetStub().CreateCompositeKey("account", []string{FromAccountId, TokenId})
// 	fmt.Printf("fromKey-->", fromKey)

// 	//Reducting the balance from From Account
// 	msg, err := contract.transactToken(ctx, fromKey, burnAmount, "-", "", "")
// 	if err != nil {
// 		//Return incase insufficient balance error or any other error while reducing the balance
// 		return msg, fmt.Errorf("%s", err.Error())
// 	}

// 	return "success", nil
// }

// //Extracting token balance of account.
// func (contract *SmartContract) BalanceOf(ctx contractapi.TransactionContextInterface, FromAccountId string, TokenId string) (int, error) {
// 	fmt.Printf("BalanceOf-->", FromAccountId, TokenId)

// 	//creating composite key AccountId+TokenId
// 	Key, _ := ctx.GetStub().CreateCompositeKey("account", []string{FromAccountId, TokenId})
// 	//Reducting the balance from From Account
// 	balance, err := contract.GetBalance(ctx, Key)
// 	if err != nil {
// 		//Return incase insufficient balance error or any other error while reducing the balance
// 		return balance, fmt.Errorf("%s", err.Error())
// 	}

// 	return balance, nil
// }

// //Function to check the Toke balance of particular account.
// func (contract *SmartContract) TotalSupply(ctx contractapi.TransactionContextInterface, TokenId string) (int, error) {

// 	fmt.Printf("GetBalance: %s", TokenId)

// 	tokenBytes, err := ctx.GetStub().GetState(TokenId)

// 	if err != nil {
// 		fmt.Errorf("Error in fecthing token supply from world state - %s", err.Error())
// 		return -1, fmt.Errorf("Error in fecthing token supply from world state - %s", err.Error())
// 	}
// 	// fmt.Printf("accountBytes: %s", accountBytes)

// 	if tokenBytes == nil {
// 		fmt.Errorf("Token does not exist with ID - %s", TokenId)
// 		return -1, nil
// 	}

// 	var token Tokens
// 	errs := json.Unmarshal(tokenBytes, &token)
// 	if errs != nil {
// 		fmt.Errorf("Error in parsing the token data - %s", errs.Error())
// 		return -1, errs
// 	}

// 	return token.TokenIssued, nil
// }

// // structure for the timelock
// type HashTimeLock struct {
// 	LockID   string `json:"lockid"`
// 	FromID   string `json:"fromid"`
// 	ToID     string `json:"toid"`
// 	TokenId  string `json:"tokenid"`
// 	Amount   int    `json:"amount"`
// 	HashLock string `json:"hashlock"`
// 	TimeLock int64  `json:"timelock"`
// 	Approved int    `json:"approved"`
// 	Minter   string `json:"minter"`
// }

// func (contract *SmartContract) TransferConditional(ctx contractapi.TransactionContextInterface, lockId string, hashlock string, timelock string, transaction string) (string, error) {

// 	var transactionData Transaction

// 	errs := json.Unmarshal([]byte(transaction), &transactionData)
// 	if errs != nil {
// 		return "fail", nil
// 	}

// 	fromKey, _ := ctx.GetStub().CreateCompositeKey("account", []string{transactionData.FromAccountId, transactionData.TokenId})
// 	toKey, _ := ctx.GetStub().CreateCompositeKey("account", []string{transactionData.ToAccountId, transactionData.TokenId})

// 	// Validate accounts
// 	accountBytes, err := ctx.GetStub().GetState(fromKey)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to read FromAccount from world state: %v", err)
// 	}
// 	if accountBytes == nil {
// 		return "", fmt.Errorf("FromAccount %s does not exist, Please create account first.", transactionData.FromAccountId)
// 	}

// 	toAccountBytes, err1 := ctx.GetStub().GetState(toKey)
// 	if err1 != nil {
// 		return "", fmt.Errorf("failed to read ToAccount from world state: %v", err1)
// 	}
// 	if toAccountBytes == nil {
// 		return "", fmt.Errorf("ToAccount %s does not exist, Please create account first.", transactionData.ToAccountId)
// 	}

// 	userId, _ := contract.getUserId(ctx)
// 	fmt.Println("TransferConditional-getUserId----", userId)

// 	isUserMinter, _ := contract.isMinter(ctx, fromKey, userId)
// 	fmt.Println("TransferConditional-isUserMinter----", isUserMinter)
// 	if isUserMinter == 1 {
// 		return "fail", fmt.Errorf("TransferConditional failed!!  User is not a minter.")
// 	}

// 	//Reducing the balance from From Account
// 	msg, err := contract.transactToken(ctx, fromKey, transactionData.Amount, "-", "", "")
// 	if err != nil {
// 		return msg, fmt.Errorf("%s", err.Error())
// 	}

// 	timeInt, err := strconv.Atoi(timelock)
// 	if err != nil {
// 		return "Error converting timeLock.", fmt.Errorf("%s", err.Error())
// 	}

// 	//calculating the expiryTime based on timeLock
// 	expiryTime := time.Now().Add(time.Minute * time.Duration(timeInt)).Unix()

// 	var hashTimeLock HashTimeLock

// 	hashTimeLock.FromID = transactionData.FromAccountId
// 	hashTimeLock.ToID = transactionData.ToAccountId
// 	hashTimeLock.Amount = transactionData.Amount
// 	hashTimeLock.TokenId = transactionData.TokenId
// 	hashTimeLock.LockID = lockId
// 	hashTimeLock.HashLock = strings.ToLower(hashlock)
// 	hashTimeLock.TimeLock = expiryTime
// 	hashTimeLock.Minter = userId
// 	hashTimeLock.Approved = 1
// 	hashTimeLockAsBytes, _ := json.Marshal(hashTimeLock)

// 	//Parking the transaction in hashTimeLock
// 	ctx.GetStub().PutState(lockId, hashTimeLockAsBytes)

// 	return "Conditional transfer successful! Hash Lock created.", nil

// }

// //Claim function to claim the hashlock by hakhathon winner
// func (contract *SmartContract) Claim(ctx contractapi.TransactionContextInterface, lock_id string, pwd string) string {

// 	hashTimeLockAsBytes, _ := ctx.GetStub().GetState(lock_id)

// 	hashTimeLock := new(HashTimeLock)
// 	_ = json.Unmarshal(hashTimeLockAsBytes, hashTimeLock)

// 	hash := sha256.Sum256([]byte(pwd))

// 	hashString := fmt.Sprintf("%x", hash)

// 	fmt.Println("Hash String:", hashString)
// 	fmt.Println("hashTimeLock---", hashTimeLock)

// 	if hashTimeLock.HashLock != hashString {

// 		fmt.Println("Invalid password:", hashString, hashTimeLock.HashLock)
// 		fmt.Println("Transaction to be reverted:")

// 		return "invalid password"
// 	}

// 	if hashTimeLock.Approved != 0 {
// 		return "Hashlock unapproved!  Tokens can not be claimed."
// 	}

// 	currTime := time.Now().Unix()

// 	fmt.Println("currTime-", currTime)

// 	//Checking if hashTimeLock expired or no
// 	if hashTimeLock.TimeLock > currTime {
// 		//Adding balance to ToAccount from hashLockTime
// 		toKey, _ := ctx.GetStub().CreateCompositeKey("account", []string{hashTimeLock.ToID, hashTimeLock.TokenId})
// 		contract.transactToken(ctx, toKey, hashTimeLock.Amount, "+", hashTimeLock.ToID, hashTimeLock.TokenId)

// 		// Deleting hashTimeLock after adding the balance to claimer account
// 		ctx.GetStub().DelState(lock_id)

// 		fmt.Println("Timelock still active. Actual transaction timestamp:", hashTimeLock.TimeLock, "Actual timelock:", currTime)
// 		return "Tokens claimed successfully!"
// 	} else {
// 		return "Tokens claimed unsuccessful!  Timelock expired."
// 	}

// }

// //Function to rever the tokens back to minters account.
// func (contract *SmartContract) Revert(ctx contractapi.TransactionContextInterface, lock_id string) (string, error) {

// 	hashTimeLockAsBytes, _ := ctx.GetStub().GetState(lock_id)

// 	hashTimeLock := new(HashTimeLock)
// 	_ = json.Unmarshal(hashTimeLockAsBytes, hashTimeLock)

// 	userId, _ := contract.getUserId(ctx)
// 	fmt.Println("TransferConditional-getUserId----", userId)

// 	isUserMinter, _ := contract.isMinter(ctx, hashTimeLock.FromID, userId)
// 	fmt.Println("TransferConditional-isUserMinter----", isUserMinter)
// 	if isUserMinter == 1 {
// 		return "fail", fmt.Errorf("Revert failed!!  User is not a minter.")
// 	}

// 	currTime := time.Now().Unix()
// 	fmt.Println("currTime-", currTime)
// 	if hashTimeLock.TimeLock > currTime {
// 		return "fail", fmt.Errorf("Timelock still not expired and token are yet open for claim.")
// 	} else {

// 		//Adding balance to FromAccount from hashLockTime
// 		fromKey, _ := ctx.GetStub().CreateCompositeKey("account", []string{hashTimeLock.FromID, hashTimeLock.TokenId})
// 		userId, _ := contract.getUserId(ctx)
// 		fmt.Println("TransferConditional-getUserId----", userId)

// 		isUserMinter, _ := contract.isMinter(ctx, fromKey, userId)
// 		fmt.Println("TransferConditional-isUserMinter----", isUserMinter)
// 		if isUserMinter == 1 {
// 			return "fail", fmt.Errorf("Revert failed!!  User is not a minter.")
// 		}

// 		contract.transactToken(ctx, fromKey, hashTimeLock.Amount, "+", "", "")

// 		// Deleting hashTimeLock after adding the balance to Minters account
// 		ctx.GetStub().DelState(lock_id)
// 		fmt.Println("Revert of Tokens successful.")

// 		return "Revert of Tokens successful to Minter.", nil
// 	}
// }
