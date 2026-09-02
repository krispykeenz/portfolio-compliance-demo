import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
export default defineConfig({plugins:[vue()],server:{proxy:{'/compliance.v1.ComplianceService':{target:'http://localhost:8080'}}},build:{commonjsOptions:{include:[/generated/,/node_modules/]}},test:{environment:'jsdom'}})
