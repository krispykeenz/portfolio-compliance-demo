import { ComplianceServiceClient } from '../generated/compliance/v1/ComplianceServiceClientPb'
import { CheckTradeRequest, GetPortfolioRequest, Money, ProposedTrade, TradeSide } from '../generated/compliance/v1/compliance_pb'

export type Holding = { id: string; name: string; value: string; exposure: string }
export type PortfolioView = { name: string; nav: string; holdings: Holding[]; remainder: string }
export type ResultView = { passed: boolean; current: string; projected: string; limit: string; explanation: string }
export interface Api { getPortfolio(): Promise<PortfolioView>; checkTrade(amountMinor: number): Promise<ResultView> }

const client = new ComplianceServiceClient(import.meta.env.VITE_GRPC_WEB_URL ?? '')

export const grpcApi: Api = {
  async getPortfolio() {
    const request = new GetPortfolioRequest(); request.setPortfolioId('portfolio-demo-001')
    const p = (await client.getPortfolio(request)).getPortfolio()
    if (!p) throw new Error('Portfolio response was empty.')
    return { name:p.getName(), nav:p.getTotalNav()?.getFormatted() ?? '', remainder:p.getRemainderDescription(), holdings:p.getHoldingsList().map(h => ({id:h.getInstrument()?.getId() ?? '', name:h.getInstrument()?.getName() ?? '', value:h.getCurrentValue()?.getFormatted() ?? '', exposure:h.getCurrentExposure()?.getFormattedPercent() ?? ''})) }
  },
  async checkTrade(amountMinor) {
    const money = new Money(); money.setMinorUnits(amountMinor); money.setCurrency('USD')
    const trade = new ProposedTrade(); trade.setPortfolioId('portfolio-demo-001'); trade.setInstrumentId('alpha'); trade.setSide(TradeSide.TRADE_SIDE_BUY); trade.setAmount(money)
    const request = new CheckTradeRequest(); request.setTrade(trade)
    const result = (await client.checkTrade(request)).getResult()
    if (!result) throw new Error('Compliance response was empty.')
    return { passed: result.getStatus() === 1, current:result.getCurrentExposure()?.getFormattedPercent() ?? '', projected:result.getProjectedExposure()?.getFormattedPercent() ?? '', limit:result.getLimitFormattedPercent(), explanation:result.getExplanation() }
  }
}
