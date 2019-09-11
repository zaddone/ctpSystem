#include <iostream>
#include <unistd.h>
#include <sys/stat.h>
//#include <stdio.h>
#include "marketspi.h"
//#include <mutex>

//mutex mut;

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
    this->mapstring["ins"] = 1;
}

void MarketSpi::routeHand(const char * data){

    cout<<"market:"<<data<<endl;
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
    case 1:{
        char *ppInstrumentID[] = {str[1]};
        this->mdApi->SubscribeMarketData(ppInstrumentID,1);
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

void MarketSpi::OnRspUserLogin(
        CThostFtdcRspUserLoginField *pRspUserLogin,
        CThostFtdcRspInfoField *pRspInfo,
        int nRequestID,
        bool bIsLast){


    cout<<"market"<<pRspInfo->ErrorID<<endl;
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
    char str[1024];
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
