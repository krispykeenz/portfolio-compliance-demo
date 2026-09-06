import { flushPromises, mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'
import App from './App.vue'
import type { Api } from './api'

const portfolio = {name:'Interview Demo Portfolio',nav:'1000000.00',remainder:'Diversified assets',holdings:[{id:'alpha',name:'Alpha Global Equity Fund',value:'90000.00',exposure:'9.00%'},{id:'beta',name:'Beta Income Fund',value:'75000.00',exposure:'7.50%'},{id:'gamma',name:'Gamma Property Fund',value:'50000.00',exposure:'5.00%'}]}
function api(result = {passed:true,current:'9.00%',projected:'10.00%',limit:'10.00%',explanation:'PASS: Alpha would move from 9.00% to 10.00%.'}): Api { return {getPortfolio:async()=>portfolio,checkTrade:async()=>result} }

describe('App', () => {
  it('renders the seeded portfolio in rands and exactly three holdings', async () => { const wrapper=mount(App,{props:{api:api()}}); await flushPromises(); expect(wrapper.text()).toContain('Interview Demo Portfolio'); expect(wrapper.text()).toContain('R1,000,000.00'); expect(wrapper.text()).toContain('Amount to add (ZAR)'); expect(wrapper.findAll('[data-holding]')).toHaveLength(3) })
  it('shows a passing response', async () => { const wrapper=mount(App,{props:{api:api()}}); await flushPromises(); await wrapper.get('form').trigger('submit'); await flushPromises(); expect(wrapper.get('[role=status]').text()).toContain('PASS') })
  it('shows a failing response', async () => { const wrapper=mount(App,{props:{api:api({passed:false,current:'9.00%',projected:'10.0001%',limit:'10.00%',explanation:'FAIL: limit exceeded.'})}}); await flushPromises(); await wrapper.get('form').trigger('submit'); await flushPromises(); expect(wrapper.get('[role=status]').text()).toContain('FAIL') })
  it('shows an API error', async () => { const broken:Api={getPortfolio:async()=>portfolio,checkTrade:async()=>{throw new Error('Service unavailable')}}; const wrapper=mount(App,{props:{api:broken}}); await flushPromises(); await wrapper.get('form').trigger('submit'); await flushPromises(); expect(wrapper.get('[role=alert]').text()).toContain('Please try again'); expect(wrapper.get('button[type=submit]').attributes('disabled')).toBeUndefined() })
  it('clears the old result when the purchase amount changes', async () => {
    const wrapper = mount(App, {props:{api:api()}})
    await flushPromises()
    await wrapper.get('form').trigger('submit')
    await flushPromises()
    expect(wrapper.find('[role=status]').exists()).toBe(true)
    await wrapper.get('input').setValue('10,001.00')
    expect(wrapper.find('[role=status]').exists()).toBe(false)
  })
  it('keeps the amount fixed while a check is pending', async () => {
    let finish!: (value: Awaited<ReturnType<Api['checkTrade']>>) => void
    const pending = new Promise<Awaited<ReturnType<Api['checkTrade']>>>(resolve => { finish = resolve })
    const wrapper = mount(App, {props:{api:{...api(),checkTrade:()=>pending}}})
    await flushPromises()
    await wrapper.get('form').trigger('submit')
    expect((wrapper.get('input').element as HTMLInputElement).matches(':disabled')).toBe(true)
    finish({passed:true,current:'9.00%',projected:'10.00%',limit:'10.00%',explanation:'Within limit.'})
    await flushPromises()
    expect((wrapper.get('input').element as HTMLInputElement).matches(':disabled')).toBe(false)
    expect(wrapper.get('[role=status]').text()).toContain('10.00%')
  })
  it('can retry loading after a connection failure', async () => {
    let attempts = 0
    const wrapper = mount(App, {props:{api:{...api(),getPortfolio:async()=>{
      if (++attempts === 1) throw new Error('Network error')
      return portfolio
    }}}})
    await flushPromises()
    expect(wrapper.find('[role=alert]').exists()).toBe(true)
    await wrapper.get('button').trigger('click')
    await flushPromises()
    expect(wrapper.findAll('[data-holding]')).toHaveLength(3)
    expect(wrapper.find('[role=alert]').exists()).toBe(false)
  })
})
