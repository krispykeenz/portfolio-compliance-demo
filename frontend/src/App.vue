<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { grpcApi, type Api, type PortfolioView, type ResultView } from './api'
import { parseMoney } from './money'
const props = withDefaults(defineProps<{api?:Api}>(), {api:()=>grpcApi})
const portfolio=ref<PortfolioView>(); const amount=ref('10,000.00'); const result=ref<ResultView>(); const error=ref(''); const loading=ref(false)
onMounted(async()=>{try{portfolio.value=await props.api.getPortfolio()}catch(e){error.value=e instanceof Error?e.message:'Unable to load portfolio.'}})
async function check(){error.value='';result.value=undefined;try{const minor=parseMoney(amount.value);loading.value=true;result.value=await props.api.checkTrade(minor)}catch(e){error.value=e instanceof Error?e.message:'Unable to check trade.'}finally{loading.value=false}}
</script>

<template>
  <main>
    <header><p class="eyebrow">PORTFOLIO COMPLIANCE</p><h1>{{portfolio?.name ?? 'Loading portfolio…'}}</h1><p v-if="portfolio" class="nav">Total NAV <strong>R{{portfolio.nav}}</strong></p></header>
    <p v-if="error" role="alert" class="error">{{error}}</p>
    <section v-if="portfolio" aria-labelledby="holdings-title"><h2 id="holdings-title">Current holdings</h2><div class="holdings"><article v-for="holding in portfolio.holdings" :key="holding.id" data-holding><h3>{{holding.name}}</h3><p>R{{holding.value}}</p><strong>{{holding.exposure}}</strong></article></div><p class="remainder">{{portfolio.remainder}}</p></section>
    <section v-if="portfolio" class="trade" aria-labelledby="trade-title"><div><p class="eyebrow">PROPOSED BUY</p><h2 id="trade-title">Alpha Global Equity Fund</h2></div><form @submit.prevent="check"><label for="amount">Trade amount (ZAR)</label><input id="amount" v-model="amount" inputmode="decimal" autocomplete="off"><div class="presets"><button v-for="preset in ['9,999.00','10,000.00','10,001.00']" :key="preset" type="button" @click="amount=preset">R{{preset}}</button></div><button class="check" type="submit" :disabled="loading">{{loading?'Checking…':'Check compliance'}}</button></form></section>
    <section v-if="result" role="status" :class="['result',result.passed?'pass':'fail']"><p class="eyebrow">RULE RESULT</p><h2>{{result.passed?'PASS':'FAIL'}}</h2><dl><div><dt>Current</dt><dd>{{result.current}}</dd></div><div><dt>Projected</dt><dd>{{result.projected}}</dd></div><div><dt>Limit</dt><dd>{{result.limit}}</dd></div></dl><p class="explanation">{{result.explanation}}</p></section>
  </main>
</template>
