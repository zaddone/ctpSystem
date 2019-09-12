#ifndef CTPSPI_H
#define CTPSPI_H
#include "marketspi.h"
#include "traderspi.h"


class ctpspi
{
public:
    ctpspi();
    ctpspi(const char *BrokerID, const char *UserID, const char * Password);
    ~ctpspi();
    void runTrader(const char *addr);
    void runMarket(const char *addr);
    void runMRecv();
    void runTRecv();
    MarketSpi * mSpi;
    TraderSpi * tSpi;

private:
    CThostFtdcReqUserLoginField UserReq;

};

#endif // CTPSPI_H
