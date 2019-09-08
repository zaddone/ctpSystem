#include <iostream>
#include "traderspi.h"
#include <unistd.h>
#include <sys/stat.h>
//#include <string.h>

using namespace std;
TraderSpi::TraderSpi(CThostFtdcReqUserLoginField *user,const char * path):socketUnixServer(path)
{

    //socketUnixServer::socketUnixServer(path);
    if (0 != access(path,0)){
        mkdir(path,0777);
    }
    this->trApi = CThostFtdcTraderApi::CreateFtdcTraderApi(path);
    this->trApi->RegisterSpi(this);
    //strcpy(&this->userReq,user);
    memset(&this->userReq,0,sizeof(this->userReq));
    strcpy(this->userReq.BrokerID,user->BrokerID);
    strcpy(this->userReq.UserID,user->UserID);
    strcpy(this->userReq.Password,user->Password);
    //memcpy(&this->userReq,user,sizeof(user));
    //this->userReq = user;

}
void TraderSpi::routeHand(const char *data){

}
int TraderSpi::getRequestID(){
    this->requestID++;
    return this->requestID;
}
void TraderSpi::OnFrontConnected(){
    cout << "Td connected"<< endl;
    int res = this->trApi->ReqUserLogin(&this->userReq,this->getRequestID());
    cout << res << endl;
}

void TraderSpi::queryInstruments()
{
    CThostFtdcQryInstrumentField req;
    memset(&req, 0, sizeof(req));
    if (!this->trApi->ReqQryInstrument(&req,this->getRequestID()))
        cout << "query Instruments error" << endl;
}

void TraderSpi::OnRspQryInstrument(
        CThostFtdcInstrumentField *pInstrument,
        CThostFtdcRspInfoField *pRspInfo,
        int nRequestID,
        bool bIsLast)
{
    //cout<< pInstrument->InstrumentName<<endl;
    //cout<< pInstrument->InstrumentID<<endl;
    char db[16] = "ins:";
    strcat(db,pInstrument->InstrumentID);
    //cout<< db <<endl;
    this->send(db);
    //cout<< db <<endl;
    //this->routeHand(db);
    //pInstrument->InstrumentID;
    //collect(pInstrument);

    //if (bIsLast)
        //signal(allInstrumentsReady);
}
void TraderSpi::OnRspUserLogin(
    CThostFtdcRspUserLoginField *pRspUserLogin,
    CThostFtdcRspInfoField *pRspInfo,
    int nRequestID,
    bool bIsLast)
{

    cout<<"trader"<<pRspInfo->ErrorID<<endl;

    char pass[]="abc2019";
    if (140==pRspInfo->ErrorID){
        CThostFtdcUserPasswordUpdateField res;
        memset(&res,0,sizeof(res));
        strcpy(res.BrokerID,this->userReq.BrokerID);
        strcpy(res.UserID,this->userReq.UserID);
        strcpy(res.OldPassword,this->userReq.Password);
        strcpy(res.NewPassword,pass);
        if (0==this->trApi->ReqUserPasswordUpdate(&res,this->getRequestID())){
            strcpy(this->userReq.Password,pass);
            this->trApi->ReqUserLogin(&this->userReq,this->getRequestID());

        }
    }else if (3 == pRspInfo->ErrorID){
        strcpy(this->userReq.Password,pass);
        this->trApi->ReqUserLogin(&this->userReq,this->getRequestID());
    }else if (0 == pRspInfo->ErrorID){
        this->queryInstruments();
    }else if (7 == pRspInfo->ErrorID){
        this->trApi->Init();
        this->trApi->ReqUserLogin(&this->userReq,this->getRequestID());
    }
    //if (0 == pRspInfo->ErrorID){
    //    this->queryInstruments();
    //};
}
