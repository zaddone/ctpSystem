#include <iostream>
#include <unistd.h>
#include <sys/stat.h>
//#include <stdio.h>
#include "marketspi.h"
#include <thread>
#include<cmath>
//#include <mutex>

//mutex mut;

MarketSpi::MarketSpi(const char * path):socketUnixServer(path){
    //if (0 != access(path,0)){
    //    mkdir(path,0777);
    //}
    this->mdApi = CThostFtdcMdApi::CreateFtdcMdApi(path,true);
    this->mdApi->RegisterSpi(this);
    memset(&this->userReq,0,sizeof(this->userReq));
    this->initMap();
}

void MarketSpi::OnRspError(CThostFtdcRspInfoField *pRspInfo, int nRequestID, bool bIsLast) {

    cout<<pRspInfo->ErrorID <<pRspInfo->ErrorMsg;
}

MarketSpi::MarketSpi(CThostFtdcReqUserLoginField *user,const char * path):socketUnixServer(path){
    if (0 != access(path,0)){
        mkdir(path,0777);
    }
    this->mdApi = CThostFtdcMdApi::CreateFtdcMdApi(path,true);
    this->mdApi->RegisterSpi(this);
    //this->userReq = user;
    memset(&this->userReq,0,sizeof(this->userReq));
    strcpy(this->userReq.BrokerID,user->BrokerID);
    strcpy(this->userReq.UserID,user->UserID);
    strcpy(this->userReq.Password,user->Password);
    //memcpy(&this->userReq,user,sizeof(user));
    this->initMap();
}
void MarketSpi::setUserReg(
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

void MarketSpi::swapPassword(){
    TThostFtdcBrokerIDType bakPass;
    strcpy(bakPass,this->userReq.Password);
    strcpy(this->userReq.Password,this->pass);
    strcpy(this->pass,bakPass);

}
void MarketSpi::Join(){
    this->mdApi->Join();
    this->send("addr");
}

void MarketSpi::run(const char *addr){
    //char _addr[1024];
    //memset(this->Addr,0,strlen(addr));
    char _addr[1024];
    strcpy(_addr,addr);
    cout<<addr<<endl;
    this->mdApi->RegisterFront(_addr);
    this->mdApi->Init();
    thread th(&MarketSpi::Join,this);
    th.detach();
    //mSpi->mdApi->Join();
}
void MarketSpi::initMap(){
    this->mapstring["ins"] = 1;
    this->mapstring["config"] = 2;
    this->mapstring["addr"] = 3;
}

void MarketSpi::routeHand(const char * data){

    cout<<"market:"<<data<<endl;
    char db[1024];
    strcpy(db,data);

    char *p;
    char sep[] = " ";
    char str[100][1024];
    p = strtok(db,sep);
    int i;
    i = 0;
    while( p != NULL ) {
        strcpy(str[i] , p);
        p = strtok(NULL, sep);
        i++;
    }
    switch (this->mapstring[str[0]]){
    case 1:{
        char *ppInstrumentID[] = {str[1]};
        this->mdApi->SubscribeMarketData(ppInstrumentID,1);
    }
        break;
    case 2:{

        cout<<"db:"<<str[1]<<str[2]<<str[3]<<str[4]<<endl;
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

int MarketSpi::getRequestID(){
    this->requestID++;
    return this->requestID;
}
void MarketSpi::OnFrontConnected(){
    cout << "Md connected"<< endl;
    int res = this->mdApi->ReqUserLogin(&this->userReq,this->getRequestID());
    cout << res << endl;
}

void MarketSpi::OnFrontDisconnected(int nReason){
    //this->send("addr");
}
void MarketSpi::OnRspUserLogin(
        CThostFtdcRspUserLoginField *pRspUserLogin,
        CThostFtdcRspInfoField *pRspInfo,
        int nRequestID,
        bool bIsLast){


    cout<<"market"<<pRspInfo->ErrorID<<endl;
    if (0==pRspInfo->ErrorID){
        this->send("ins");
    }
}

void MarketSpi::OnRspUnSubMarketData(
        CThostFtdcSpecificInstrumentField *pSpecificInstrument,
        CThostFtdcRspInfoField *pRspInfo,
        int nRequestID,
        bool bIsLast){


}
void MarketSpi::OnRspSubMarketData(
        CThostFtdcSpecificInstrumentField *pSpecificInstrument,
        CThostFtdcRspInfoField *pRspInfo,
        int nRequestID,
        bool bIsLast){
    printf("%s %s\n",pSpecificInstrument->InstrumentID,pRspInfo->ErrorMsg);

}

void MarketSpi::OnRtnDepthMarketData(CThostFtdcDepthMarketDataField *pDepthMarketData){
    //printf("%s %s %d %lf %lf\n",
    //       pDepthMarketData->InstrumentID,
    //       pDepthMarketData->UpdateTime,
    //       pDepthMarketData->UpdateMillisec,
    //       pDepthMarketData->AskPrice1,
    //       pDepthMarketData->BidPrice1
    //       );

    char str[8192];
    sprintf(str,"market %s,%sT%s,%lf,%lf",
           pDepthMarketData->InstrumentID,
           pDepthMarketData->TradingDay,
           pDepthMarketData->UpdateTime,
           pDepthMarketData->AskPrice1,
           pDepthMarketData->BidPrice1
            );
    //cout<<str<<endl;
    this->send(str);
}
