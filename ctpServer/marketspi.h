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
    //MarketSpi(CThostFtdcReqUserLoginField *user, const char *path);
    CThostFtdcMdApi *mdApi;
    virtual void routeHand(const char *data);
    int getRequestID();
    void run(const char *addr);
    virtual void OnFrontDisconnected(int nReason);
    virtual void OnFrontConnected();
    virtual void OnRspUserLogin(
        CThostFtdcRspUserLoginField *pRspUserLogin,
        CThostFtdcRspInfoField *pRspInfo,
        int nRequestID,
        bool bIsLast);
    virtual void OnRspSubMarketData(
            CThostFtdcSpecificInstrumentField *pSpecificInstrument,
            CThostFtdcRspInfoField *pRspInfo,
            int nRequestID,
            bool bIsLast) ;
    virtual void OnRspUnSubMarketData(
            CThostFtdcSpecificInstrumentField *pSpecificInstrument,
            CThostFtdcRspInfoField *pRspInfo,
            int nRequestID,
            bool bIsLast);
    virtual void OnRtnDepthMarketData(CThostFtdcDepthMarketDataField *pDepthMarketData);
    virtual void OnRspError(CThostFtdcRspInfoField *pRspInfo, int nRequestID, bool bIsLast) ;
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
    void setUserReg(
        const char * brokerID,
        const char * userID,
        const char * password,
        const char * passwordBak);
    const char * path;
    void stop();

};

#endif // MARKETSPI_H
