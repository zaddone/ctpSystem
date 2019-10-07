#ifndef MARKETSPI_H
#define MARKETSPI_H
#include "ThostFtdcMdApi.h"
#include "socketunixserver.h"
#include <map>
using namespace std;

class MarketSpi : public socketUnixServer,CThostFtdcMdSpi
{
public:
    MarketSpi(const char *path);
    MarketSpi(const char *brokerID,
              const char *userID,
              const char *password,
              const char *passwordBak,
              const char *addr,
              const char *path);
    CThostFtdcMdApi *mdApi;
    virtual void routeHand(const char *data);
    int getRequestID();
    virtual void OnFrontConnected();
    virtual void OnRspUserLogin(
        CThostFtdcRspUserLoginField *pRspUserLogin,
        CThostFtdcRspInfoField *pRspInfo,
        int nRequestID,
        bool bIsLast);
    //virtual void OnRspSubMarketData(
    //        CThostFtdcSpecificInstrumentField *pSpecificInstrument,
    //        CThostFtdcRspInfoField *pRspInfo,
    //        int nRequestID,
    //        bool bIsLast) ;
    //virtual void OnRspUnSubMarketData(
    //        CThostFtdcSpecificInstrumentField *pSpecificInstrument,
    //        CThostFtdcRspInfoField *pRspInfo,
    //        int nRequestID,
    //        bool bIsLast);
    virtual void OnRtnDepthMarketData(CThostFtdcDepthMarketDataField *pDepthMarketData);
    virtual void OnRspError(CThostFtdcRspInfoField *pRspInfo, int nRequestID, bool bIsLast) ;
    virtual void OnFrontDisconnected(int nReason);

    //this->stop();
    //virtual const char *GetTradingDay();

private:
    //char * Addr;
    void swapPassword();
    TThostFtdcPasswordType pass;
    int requestID;
    CThostFtdcReqUserLoginField userReq;
    map<string , int >mapstring;
    void initMap();
    void Join();
    void run();
    void setUserReg(
        const char * brokerID,
        const char * userID,
        const char * password,
        const char * passwordBak);
    const char * path;
    void stop();
    void reqUserLogin();
    //char TradingDay[8];
    //bool Login;
    const char * Addr;
    void subscribeMarketData(char *ins);
    void help();

};

#endif // MARKETSPI_H
