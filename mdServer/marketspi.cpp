#include <iostream>
#include <unistd.h>
//#include <sys/stat.h>
//#include <stdio.h>
#include <thread>
#include "marketspi.h"
//#include<cmath>
//#include <mutex>

//mutex mut;

bool IsFlowControl(int iResult)
{
    return ((iResult == -2) || (iResult == -3));
}
MarketSpi::MarketSpi(const char * path):socketUnixServer(path){
    //if (0 != access(path,0)){
    //    mkdir(path,0777);
    //}
    this->path = path;

    memset(&this->userReq,0,sizeof(this->userReq));
    this->initMap();
}

void MarketSpi::OnRspError(CThostFtdcRspInfoField *pRspInfo, int nRequestID, bool bIsLast) {

    cout<<pRspInfo->ErrorID <<pRspInfo->ErrorMsg;
}

MarketSpi::MarketSpi(
        const char * brokerID,
        const char * userID,
        const char *password,
        const char *passwordBak,
        const char *addr,
        const char * path):socketUnixServer(path){
    this->mdApi = NULL;
    //this->TradingDay = NULL;
    //this->Login = false;
    this->path = path;
    this->Addr = addr;
    this->initMap();
    memset(&this->userReq,0,sizeof(this->userReq));
    this->setUserReg(brokerID,userID,password,passwordBak);
    this->run();
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
    //this->stop();
    //this->send("addr");
    cout<<"stop market"<<endl;
}
void MarketSpi::stop(){
    if (this->mdApi==NULL)return;
    this->mdApi->RegisterSpi(NULL);
    this->mdApi->Release();
    this->mdApi = NULL;
    //this->over =  true;
}

void MarketSpi::run(){

    if (this->mdApi != NULL) return;

    this->mdApi = CThostFtdcMdApi::CreateFtdcMdApi(this->path,true);
    this->mdApi->RegisterSpi(this);
    char _addr[1024];
    strcpy(_addr,Addr);
    //cout<<addr<<endl;
    this->mdApi->RegisterFront(_addr);
    this->mdApi->Init();
    thread th(&MarketSpi::Join,this);
    th.detach();
    //mSpi->mdApi->Join();
}
void MarketSpi::initMap(){
    this->mapstring["help"] = 999;
    this->mapstring["stop"] = 100;
    this->mapstring["ins"] = 1;
    this->mapstring["config"] = 2;
    this->mapstring["addr"] = 3;
}
void MarketSpi::help(){
    map<string , int>::iterator iter;
    for(iter = mapstring.begin(); iter != mapstring.end(); iter++)
          cout<<"help:"<<iter->first<<endl;

}

void MarketSpi::routeHand(const char * data){

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
    case 999:{
        this->help();
        break;
    }
    case 100:{
        this->stop();
        break;
    }
    case 1:{
        this->subscribeMarketData(str[1]);
        break;
    }
    case 2:{
        //cout<<"db:"<<str[1]<<str[2]<<str[3]<<str[4]<<endl;
        this->setUserReg(str[1],str[2],str[3],str[4]);
        this->Addr = str[5];
        //this->stop();
        this->run();
        break;
    }
    case 3:{
        //this->setUserReg(str[1],str[2],str[3],str[4]);
        this->Addr = str[1];
        //this->stop();
        this->run();
        break;
    }
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
    cout<<"conn"<<endl;
    this->reqUserLogin();
}

void MarketSpi::reqUserLogin(){
    while (true)
    {
        int iResult = this->mdApi->ReqUserLogin(&this->userReq,this->getRequestID());
        if (!IsFlowControl(iResult))
        {
            break;
        }
        else
        {
            sleep(1);
        }
    }
}

void MarketSpi::OnFrontDisconnected(int nReason){
    cout<<"disconnected:"<<nReason<<endl;
    this->stop();
}

void MarketSpi::OnRspUserLogin(
        CThostFtdcRspUserLoginField *pRspUserLogin,
        CThostFtdcRspInfoField *pRspInfo,
        int nRequestID,
        bool bIsLast){
    cout<<"mk "<<pRspInfo->ErrorID<<endl;
    if (pRspInfo && 0!=pRspInfo->ErrorID){
        cout<<pRspInfo->ErrorMsg<<endl;
        //this->run();
        return;
    }

    char trading[20]="TDay ";
    strcat(trading,this->mdApi->GetTradingDay());
    this->send(trading);
    //strcpy(this->TradingDay,this->mdApi->GetTradingDay());
    //cout<<"ins"<<endl;

}

void MarketSpi::OnRtnDepthMarketData(CThostFtdcDepthMarketDataField *pDepthMarketData){

    char str[8192];
    sprintf(str,"market %s,%sT%s,%lf,%lf,%lf,%lf",
           pDepthMarketData->InstrumentID,
           pDepthMarketData->TradingDay,
           pDepthMarketData->UpdateTime,
           pDepthMarketData->AskPrice1,
           pDepthMarketData->BidPrice1,
           pDepthMarketData->UpperLimitPrice,
           pDepthMarketData->LowerLimitPrice
    );
    this->send(str);

}

void MarketSpi::subscribeMarketData(char * ins){

    char *ppInstrumentID[] = {ins};
    while (true)
    {
        //cout<<ins<<endl;
        int iResult = this->mdApi->SubscribeMarketData(ppInstrumentID,1);
        if (!IsFlowControl(iResult))
            break;
        else
            sleep(1);
    }

}
