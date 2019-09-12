#include <iostream>
#include "traderspi.h"
#include <unistd.h>
#include <sys/stat.h>
//#include <string.h>

using namespace std;
TraderSpi::TraderSpi(CThostFtdcReqUserLoginField *user,const char * path):socketUnixServer(path)
{

    //socketUnixServer::socketUnixServer(path);
    //if (0 != access(path,0)){
    // 	  mkdir(path,0777);
    //}
    this->trApi = CThostFtdcTraderApi::CreateFtdcTraderApi(path);
    this->trApi->RegisterSpi(this);
    //strcpy(&this->userReq,user);
    memset(&this->userReq,0,sizeof(this->userReq));
    strcpy(this->userReq.BrokerID,user->BrokerID);
    strcpy(this->userReq.UserID,user->UserID);
    strcpy(this->userReq.Password,user->Password);
    strcpy(this->pass , "abc2019");
    //memcpy(&this->userReq,user,sizeof(user));
    //this->userReq = user;

}
void TraderSpi::routeHand(const char *data){

    cout<<"trader:"<<data<<endl;
    char db[1024];
    strcpy(db,data);
    cout<<"db:"<<db<<endl;
    char *p;
    char sep[] = " ";
    char str[2][1024];
    p = strtok(db,sep);
    int i;
    i = 0;
    while( p != NULL ) {
        strcpy(str[i] , p);
        p = strtok(NULL, sep);
        i++;
    }
    switch (this->mapstring[str[0]]){

    case 2:{
        this->setUserReg(str[1],str[2],str[3],str[4]);
        this->run(str[5]);
    }
        break;
    case 3:{
        //this->setUserReg(str[1],str[2],str[3],str[4]);
        this->run(str[1]);
    }
        break;
    default:
        printf("default %s",data);
        break;
    }


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
    if (this->queryIns)return;
    this->queryIns = true;
    CThostFtdcQryInstrumentField req;
    memset(&req, 0, sizeof(req));
    this->trApi->ReqQryInstrument(&req,this->getRequestID());
}

void TraderSpi::OnRspQryInstrument(
        CThostFtdcInstrumentField *pInstrument,
        CThostFtdcRspInfoField *pRspInfo,
        int nRequestID,
        bool bIsLast)
{
    //cout<< pInstrument->InstrumentName<<endl;
    //cout<< pInstrument->InstrumentID<<endl;
    char db[16] = "ins ";
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
void TraderSpi::swapPassword(){
    TThostFtdcBrokerIDType bakPass;
    strcpy(bakPass,this->userReq.Password);
    strcpy(this->userReq.Password,this->pass);
    strcpy(this->pass,bakPass);

}
void TraderSpi::OnRspUserLogin(
    CThostFtdcRspUserLoginField *pRspUserLogin,
    CThostFtdcRspInfoField *pRspInfo,
    int nRequestID,
    bool bIsLast)
{

    cout<<"trader"<<pRspInfo->ErrorID<<endl;

    //char pass[]="abc2019";
    if (140==pRspInfo->ErrorID){
        CThostFtdcUserPasswordUpdateField res;
        memset(&res,0,sizeof(res));
        strcpy(res.BrokerID,this->userReq.BrokerID);
        strcpy(res.UserID,this->userReq.UserID);
        strcpy(res.OldPassword,this->userReq.Password);
        strcpy(res.NewPassword,pass);
        if (0==this->trApi->ReqUserPasswordUpdate(&res,this->getRequestID())){
            this->swapPassword();
            //TThostFtdcBrokerIDType bakPass;
            //strcpy(bakPass,this->userReq.Password);
            //strcpy(this->userReq.Password,pass);
            //strcpy(pass,bakPass);
            this->trApi->ReqUserLogin(&this->userReq,this->getRequestID());

        }
    }else if (3 == pRspInfo->ErrorID){
        strcpy(this->userReq.Password,pass);
        this->trApi->ReqUserLogin(&this->userReq,this->getRequestID());
    }else if (0 == pRspInfo->ErrorID){
        this->queryInstruments();
    }else if (7 == pRspInfo->ErrorID){
        this->swapPassword();
        //TThostFtdcBrokerIDType bakPass;
        //strcpy(bakPass,this->userReq.Password);
        //strcpy(this->userReq.Password,pass);
        //strcpy(pass,bakPass);
        this->trApi->ReqUserLogin(&this->userReq,this->getRequestID());
    }
    //if (0 == pRspInfo->ErrorID){
    //    this->queryInstruments();
    //};
}

void TraderSpi::run(const char *addr){
    char _addr[1024];
    strcpy(_addr,addr);
    cout<<addr<<endl;
    this->trApi->RegisterFront(_addr);
    this->trApi->Init();
    //mSpi->mdApi->Join();
}
void TraderSpi::initMap(){
    //this->mapstring["ins"] = 1;
    this->mapstring["config"] = 2;
    this->mapstring["addr"] = 3;
}
void TraderSpi::setUserReg(
        const char * brokerID,
        const char * userID,
        const char *password,
        const char *passwordBak){
    //memset(&this->userReq,0,sizeof(this->userReq));
    strcpy(this->userReq.BrokerID,brokerID);
    strcpy(this->userReq.UserID,userID);
    strcpy(this->userReq.Password,password);
    strcpy(this->pass,passwordBak);
}
