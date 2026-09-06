<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { grpcApi, type Api, type PortfolioView, type ResultView } from './api'
import { formatMoney, parseMoney } from './money'

const props = withDefaults(defineProps<{ api?: Api }>(), { api: () => grpcApi })
const portfolio = ref<PortfolioView>()
const amount = ref('10,000.00')
const result = ref<ResultView>()
const loading = ref(false)
const loadError = ref(false)
const error = ref('')
const amountError = ref('')

async function loadPortfolio() {
  loadError.value = false
  try { portfolio.value = await props.api.getPortfolio() }
  catch { loadError.value = true }
}
onMounted(loadPortfolio)
watch(amount, () => { result.value = undefined; error.value = ''; amountError.value = '' })

async function check() {
  if (loading.value) return
  error.value = ''; amountError.value = ''; result.value = undefined
  let minor: number
  try { minor = parseMoney(amount.value) }
  catch (e) { amountError.value = (e as Error).message; return }
  loading.value = true
  try { result.value = await props.api.checkTrade(minor) }
  catch { error.value = 'We couldn’t complete your check. Please try again.' }
  finally { loading.value = false }
}
</script>

<template>
  <a class="skip-link" href="#main">Skip to content</a>
  <header class="topbar">
    <div class="topbar-inner">
      <div class="brand">
        <svg class="brand-mark" viewBox="0 0 32 32" fill="none" aria-hidden="true"><path d="M5 24V14h5v10M14 24V8h5v16M23 24V3h5v21" stroke="currentColor" stroke-width="3"/></svg>
        <span>Portfolio<span class="brand-caption">INVESTMENT CHECK</span></span>
      </div>
      <span class="demo-badge"><span aria-hidden="true"></span>Demo portfolio</span>
    </div>
  </header>

  <main id="main">
    <div class="page-heading">
      <div><p class="eyebrow">A LITTLE CLARITY BEFORE YOU INVEST</p><h1>Portfolio overview</h1><p class="intro">Explore a sample portfolio and see how a purchase changes its balance.</p></div>
      <a class="text-link" href="#how-it-works">How this works <span aria-hidden="true">↗</span></a>
    </div>

    <div v-if="loadError" class="load-state panel" role="alert">
      <h2>We couldn’t load the portfolio.</h2><p>Please check your connection and try again.</p>
      <button type="button" class="primary-button" @click="loadPortfolio">Try again</button>
    </div>
    <div v-else-if="!portfolio" class="load-state panel" role="status"><span class="spinner" aria-hidden="true"></span><p>Loading your sample portfolio…</p></div>

    <template v-else>
      <section class="overview" aria-label="Portfolio summary">
        <div class="portfolio-value">
          <p class="eyebrow">TOTAL INVESTMENT VALUE</p>
          <p class="balance">{{ formatMoney(portfolio.nav) }}</p>
          <p class="portfolio-name">{{ portfolio.name }}</p>
          <div class="allocation" aria-hidden="true"><span v-for="(holding, index) in portfolio.holdings" :key="holding.id" :class="`fund-color-${index}`" :style="{ width: holding.exposure }"></span></div>
          <p class="allocation-caption">All your investments, added together.</p>
        </div>
        <div class="portfolio-rule">
          <p class="eyebrow">THE RULE WE’RE CHECKING</p>
          <p class="rule-number">10<span>%</span></p>
          <h2>Maximum in one fund</h2>
          <p>No single fund should make up more than 10% of the total in this demo. Exactly 10% is allowed.</p>
        </div>
      </section>

      <div class="workspace">
        <div class="portfolio-column">
          <section class="panel investments" aria-labelledby="investments-title">
            <div class="section-heading"><div><p class="eyebrow">WHERE THE MONEY IS</p><h2 id="investments-title">Your investments</h2></div><span class="subtle-tag">ZAR</span></div>
            <p class="section-description">A portfolio is the collection of investments you own.</p>
            <table>
              <caption class="sr-only">Current fund values and their share of all investments</caption>
              <thead><tr><th scope="col">Investment fund</th><th scope="col">Value</th><th scope="col">Share of total</th></tr></thead>
              <tbody><tr v-for="(holding, index) in portfolio.holdings" :key="holding.id" data-holding>
                <th scope="row"><span class="fund-name"><span :class="['fund-dot', `fund-color-${index}`]" aria-hidden="true"></span>{{ holding.name }}</span></th>
                <td>{{ formatMoney(holding.value) }}</td><td><span class="share">{{ holding.exposure }}</span></td>
              </tr></tbody>
            </table>
            <div class="table-note"><span class="other-dot" aria-hidden="true"></span><p>The rest is spread across other investments that aren’t listed here.</p></div>
          </section>

          <section id="how-it-works" class="guide" aria-labelledby="guide-title">
            <p class="eyebrow">THE IDEA, IN PLAIN ENGLISH</p><h2 id="guide-title">A balance check for your money.</h2>
            <p>A fund pools money into a group of investments. This check tells you whether a purchase would put more than the allowed share of your money into one fund.</p>
            <ol>
              <li><span>01</span><div><h3>Choose an amount</h3><p>Try adding money to Alpha Global Equity Fund.</p></div></li>
              <li><span>02</span><div><h3>See what would change</h3><p>The purchase uses money from the rest of the portfolio, so the total stays the same.</p></div></li>
              <li><span>03</span><div><h3>Check it against the limit</h3><p>We compare the fund’s new share with the 10% limit. No purchase is made.</p></div></li>
            </ol>
            <p class="guide-note">This checks one rule. It doesn’t predict returns or tell you whether an investment is a good choice.</p>
          </section>
        </div>

        <section class="panel trade" aria-labelledby="trade-title">
          <div class="section-heading"><div><p class="eyebrow">EXPLORE A WHAT-IF</p><h2 id="trade-title">Try a purchase</h2></div><span class="simulation-icon" aria-hidden="true">↗</span></div>
          <p class="section-description">See whether adding to this fund stays within the limit.</p>
          <div class="selected-fund"><span class="fund-dot fund-color-0" aria-hidden="true"></span><div><span class="field-caption">ADDING TO</span><strong>Alpha Global Equity Fund</strong></div></div>
          <form @submit.prevent="check" :aria-busy="loading">
            <fieldset :disabled="loading">
              <legend class="sr-only">Proposed purchase amount</legend>
              <label for="amount">Amount to add (ZAR)</label>
              <div :class="['amount-input', { invalid: amountError }]"><span aria-hidden="true">R</span><input id="amount" v-model="amount" inputmode="decimal" autocomplete="off" :aria-invalid="amountError ? true : undefined" :aria-describedby="amountError ? 'amount-help amount-error' : 'amount-help'"></div>
              <p id="amount-help" class="input-help">Enter an amount in South African rand.</p>
              <p v-if="amountError" id="amount-error" class="error" role="alert">{{ amountError }}</p>
              <p class="preset-label">Or try an example</p>
              <div class="presets"><button v-for="preset in ['9,999.00', '10,000.00', '10,001.00']" :key="preset" type="button" :aria-pressed="amount === preset" @click="amount = preset">R{{ preset }}</button></div>
              <button class="primary-button check" type="submit" :disabled="loading"><span v-if="loading" class="spinner" aria-hidden="true"></span>{{ loading ? 'Checking your purchase…' : 'Check this purchase' }}<span v-if="!loading" aria-hidden="true">→</span></button>
            </fieldset>
          </form>
          <p class="simulation-note"><svg viewBox="0 0 16 16" fill="none" aria-hidden="true"><rect x="3.5" y="7" width="9" height="7" rx="1.5" stroke="currentColor"/><path d="M5.5 7V4.5a2.5 2.5 0 0 1 5 0V7" stroke="currentColor"/></svg>Simulation only. Your investments won’t change.</p>
          <p v-if="error" role="alert" class="error">{{ error }}</p>

          <section v-if="result" role="status" :class="['result', result.passed ? 'pass' : 'fail']" aria-labelledby="result-title">
            <div class="result-heading"><span class="result-symbol" aria-hidden="true">{{ result.passed ? '✓' : '!' }}</span><div><p class="eyebrow">{{ result.passed ? 'PASS' : 'FAIL' }} · PURCHASE CHECK</p><h3 id="result-title">{{ result.passed ? 'Within the limit' : 'Over the limit' }}</h3></div></div>
            <dl><div><dt>Share now</dt><dd>{{ result.current }}</dd></div><div><dt>After purchase</dt><dd>{{ result.projected }}</dd></div><div><dt>Maximum</dt><dd>{{ result.limit }}</dd></div></dl>
            <p class="explanation">{{ result.explanation }}</p>
          </section>
          <div v-else class="result-placeholder"><span aria-hidden="true">◎</span><p>{{ loading ? 'Comparing the new share with the limit…' : 'Your result will explain how the fund’s share would change.' }}</p></div>
        </section>
      </div>
      <footer class="page-footer"><span>Sample data · South African rand (ZAR)</span><span>Every check starts with the same portfolio.</span></footer>
    </template>
  </main>
</template>
