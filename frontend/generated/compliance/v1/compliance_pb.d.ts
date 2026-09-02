import * as jspb from 'google-protobuf'



export class GetPortfolioRequest extends jspb.Message {
  getPortfolioId(): string;
  setPortfolioId(value: string): GetPortfolioRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetPortfolioRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetPortfolioRequest): GetPortfolioRequest.AsObject;
  static serializeBinaryToWriter(message: GetPortfolioRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetPortfolioRequest;
  static deserializeBinaryFromReader(message: GetPortfolioRequest, reader: jspb.BinaryReader): GetPortfolioRequest;
}

export namespace GetPortfolioRequest {
  export type AsObject = {
    portfolioId: string;
  };
}

export class Money extends jspb.Message {
  getMinorUnits(): number;
  setMinorUnits(value: number): Money;

  getCurrency(): string;
  setCurrency(value: string): Money;

  getFormatted(): string;
  setFormatted(value: string): Money;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Money.AsObject;
  static toObject(includeInstance: boolean, msg: Money): Money.AsObject;
  static serializeBinaryToWriter(message: Money, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Money;
  static deserializeBinaryFromReader(message: Money, reader: jspb.BinaryReader): Money;
}

export namespace Money {
  export type AsObject = {
    minorUnits: number;
    currency: string;
    formatted: string;
  };
}

export class Exposure extends jspb.Message {
  getInstrumentValueMinorUnits(): number;
  setInstrumentValueMinorUnits(value: number): Exposure;

  getPortfolioNavMinorUnits(): number;
  setPortfolioNavMinorUnits(value: number): Exposure;

  getFormattedPercent(): string;
  setFormattedPercent(value: string): Exposure;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Exposure.AsObject;
  static toObject(includeInstance: boolean, msg: Exposure): Exposure.AsObject;
  static serializeBinaryToWriter(message: Exposure, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Exposure;
  static deserializeBinaryFromReader(message: Exposure, reader: jspb.BinaryReader): Exposure;
}

export namespace Exposure {
  export type AsObject = {
    instrumentValueMinorUnits: number;
    portfolioNavMinorUnits: number;
    formattedPercent: string;
  };
}

export class Instrument extends jspb.Message {
  getId(): string;
  setId(value: string): Instrument;

  getName(): string;
  setName(value: string): Instrument;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Instrument.AsObject;
  static toObject(includeInstance: boolean, msg: Instrument): Instrument.AsObject;
  static serializeBinaryToWriter(message: Instrument, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Instrument;
  static deserializeBinaryFromReader(message: Instrument, reader: jspb.BinaryReader): Instrument;
}

export namespace Instrument {
  export type AsObject = {
    id: string;
    name: string;
  };
}

export class Holding extends jspb.Message {
  getInstrument(): Instrument | undefined;
  setInstrument(value?: Instrument): Holding;
  hasInstrument(): boolean;
  clearInstrument(): Holding;

  getCurrentValue(): Money | undefined;
  setCurrentValue(value?: Money): Holding;
  hasCurrentValue(): boolean;
  clearCurrentValue(): Holding;

  getCurrentExposure(): Exposure | undefined;
  setCurrentExposure(value?: Exposure): Holding;
  hasCurrentExposure(): boolean;
  clearCurrentExposure(): Holding;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Holding.AsObject;
  static toObject(includeInstance: boolean, msg: Holding): Holding.AsObject;
  static serializeBinaryToWriter(message: Holding, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Holding;
  static deserializeBinaryFromReader(message: Holding, reader: jspb.BinaryReader): Holding;
}

export namespace Holding {
  export type AsObject = {
    instrument?: Instrument.AsObject;
    currentValue?: Money.AsObject;
    currentExposure?: Exposure.AsObject;
  };
}

export class Portfolio extends jspb.Message {
  getId(): string;
  setId(value: string): Portfolio;

  getName(): string;
  setName(value: string): Portfolio;

  getTotalNav(): Money | undefined;
  setTotalNav(value?: Money): Portfolio;
  hasTotalNav(): boolean;
  clearTotalNav(): Portfolio;

  getHoldingsList(): Array<Holding>;
  setHoldingsList(value: Array<Holding>): Portfolio;
  clearHoldingsList(): Portfolio;
  addHoldings(value?: Holding, index?: number): Holding;

  getRemainderDescription(): string;
  setRemainderDescription(value: string): Portfolio;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Portfolio.AsObject;
  static toObject(includeInstance: boolean, msg: Portfolio): Portfolio.AsObject;
  static serializeBinaryToWriter(message: Portfolio, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Portfolio;
  static deserializeBinaryFromReader(message: Portfolio, reader: jspb.BinaryReader): Portfolio;
}

export namespace Portfolio {
  export type AsObject = {
    id: string;
    name: string;
    totalNav?: Money.AsObject;
    holdingsList: Array<Holding.AsObject>;
    remainderDescription: string;
  };
}

export class GetPortfolioResponse extends jspb.Message {
  getPortfolio(): Portfolio | undefined;
  setPortfolio(value?: Portfolio): GetPortfolioResponse;
  hasPortfolio(): boolean;
  clearPortfolio(): GetPortfolioResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetPortfolioResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetPortfolioResponse): GetPortfolioResponse.AsObject;
  static serializeBinaryToWriter(message: GetPortfolioResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetPortfolioResponse;
  static deserializeBinaryFromReader(message: GetPortfolioResponse, reader: jspb.BinaryReader): GetPortfolioResponse;
}

export namespace GetPortfolioResponse {
  export type AsObject = {
    portfolio?: Portfolio.AsObject;
  };
}

export class ProposedTrade extends jspb.Message {
  getPortfolioId(): string;
  setPortfolioId(value: string): ProposedTrade;

  getInstrumentId(): string;
  setInstrumentId(value: string): ProposedTrade;

  getSide(): TradeSide;
  setSide(value: TradeSide): ProposedTrade;

  getAmount(): Money | undefined;
  setAmount(value?: Money): ProposedTrade;
  hasAmount(): boolean;
  clearAmount(): ProposedTrade;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProposedTrade.AsObject;
  static toObject(includeInstance: boolean, msg: ProposedTrade): ProposedTrade.AsObject;
  static serializeBinaryToWriter(message: ProposedTrade, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProposedTrade;
  static deserializeBinaryFromReader(message: ProposedTrade, reader: jspb.BinaryReader): ProposedTrade;
}

export namespace ProposedTrade {
  export type AsObject = {
    portfolioId: string;
    instrumentId: string;
    side: TradeSide;
    amount?: Money.AsObject;
  };
}

export class CheckTradeRequest extends jspb.Message {
  getTrade(): ProposedTrade | undefined;
  setTrade(value?: ProposedTrade): CheckTradeRequest;
  hasTrade(): boolean;
  clearTrade(): CheckTradeRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CheckTradeRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CheckTradeRequest): CheckTradeRequest.AsObject;
  static serializeBinaryToWriter(message: CheckTradeRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CheckTradeRequest;
  static deserializeBinaryFromReader(message: CheckTradeRequest, reader: jspb.BinaryReader): CheckTradeRequest;
}

export namespace CheckTradeRequest {
  export type AsObject = {
    trade?: ProposedTrade.AsObject;
  };
}

export class RuleResult extends jspb.Message {
  getRuleId(): string;
  setRuleId(value: string): RuleResult;

  getStatus(): ComplianceStatus;
  setStatus(value: ComplianceStatus): RuleResult;

  getReasonCode(): ReasonCode;
  setReasonCode(value: ReasonCode): RuleResult;

  getCurrentExposure(): Exposure | undefined;
  setCurrentExposure(value?: Exposure): RuleResult;
  hasCurrentExposure(): boolean;
  clearCurrentExposure(): RuleResult;

  getProjectedExposure(): Exposure | undefined;
  setProjectedExposure(value?: Exposure): RuleResult;
  hasProjectedExposure(): boolean;
  clearProjectedExposure(): RuleResult;

  getLimitBasisPoints(): number;
  setLimitBasisPoints(value: number): RuleResult;

  getLimitFormattedPercent(): string;
  setLimitFormattedPercent(value: string): RuleResult;

  getExplanation(): string;
  setExplanation(value: string): RuleResult;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RuleResult.AsObject;
  static toObject(includeInstance: boolean, msg: RuleResult): RuleResult.AsObject;
  static serializeBinaryToWriter(message: RuleResult, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RuleResult;
  static deserializeBinaryFromReader(message: RuleResult, reader: jspb.BinaryReader): RuleResult;
}

export namespace RuleResult {
  export type AsObject = {
    ruleId: string;
    status: ComplianceStatus;
    reasonCode: ReasonCode;
    currentExposure?: Exposure.AsObject;
    projectedExposure?: Exposure.AsObject;
    limitBasisPoints: number;
    limitFormattedPercent: string;
    explanation: string;
  };
}

export class CheckTradeResponse extends jspb.Message {
  getResult(): RuleResult | undefined;
  setResult(value?: RuleResult): CheckTradeResponse;
  hasResult(): boolean;
  clearResult(): CheckTradeResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CheckTradeResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CheckTradeResponse): CheckTradeResponse.AsObject;
  static serializeBinaryToWriter(message: CheckTradeResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CheckTradeResponse;
  static deserializeBinaryFromReader(message: CheckTradeResponse, reader: jspb.BinaryReader): CheckTradeResponse;
}

export namespace CheckTradeResponse {
  export type AsObject = {
    result?: RuleResult.AsObject;
  };
}

export enum TradeSide {
  TRADE_SIDE_UNSPECIFIED = 0,
  TRADE_SIDE_BUY = 1,
  TRADE_SIDE_SELL = 2,
}
export enum ComplianceStatus {
  COMPLIANCE_STATUS_UNSPECIFIED = 0,
  COMPLIANCE_STATUS_PASS = 1,
  COMPLIANCE_STATUS_FAIL = 2,
}
export enum ReasonCode {
  REASON_CODE_UNSPECIFIED = 0,
  REASON_CODE_WITHIN_LIMIT = 1,
  REASON_CODE_LIMIT_EXCEEDED = 2,
}
