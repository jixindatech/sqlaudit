<template>
  <div class="dashboard-container">
    <panel-group
      :event-total="eventTotal"
      :allowed-total="allowedTotal"
      :denied-total="deniedTotal"
      :unknown-total="unknownTotal"
      :options="options"
      :query-data="queryData"
    />
    <el-row :gutter="40">
      <el-col :xs="24" :sm="24" :lg="12">
        <el-card>
          <IPBarChart v-if="flag" :data="ip" />
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="24" :lg="12">
        <el-card>
          <UserBarChart v-if="flag" :data="user" />
        </el-card>
      </el-col>
    </el-row>
    <!--
    <el-row :gutter="40" style="margin-top:30px">
      <el-col :xs="24" :sm="24" :lg="12">
        <el-card>
          <bar-chart v-if="flag" />
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="24" :lg="12">
        <el-card>
          <pie-chart v-if="flag" :pie-data="pieData" />
        </el-card>
      </el-col>
    </el-row>
    -->
    <!--<el-row style="background:#fff;top-bittom:120px; padding:16px 16px 0;margin-bottom:32px;">-->
    <el-row style="margin-top:30px">
      <el-card>
        <OPLineChart :data="opinfo" />
      </el-card>
    </el-row>
    <el-row style="margin-top:30px">
      <el-card>
        <TypeLineChart :data="typeinfo" />
      </el-card>
    </el-row>
  </div>
</template>

<script>
import PanelGroup from './components/PanelGroup'
// import PieChart from './components/PieChart'
// import BarChart from './components/BarChart'
// import LineChart from './components/LineChart.vue'
import IPBarChart from './components/IPBarChart.vue'
import UserBarChart from './components/UserBarChart.vue'
import OPLineChart from './components/OPLineChart.vue'
import TypeLineChart from './components/TypeLineChart.vue'

// import host from '@/api/host'
import api from '@/api/statics'
import { sqlStrOpMap } from '@/utils/const'

export default {
  name: 'Dashboard',
  components: { PanelGroup, /* PieChart, BarChart, LineChart,*/ IPBarChart, UserBarChart, OPLineChart, TypeLineChart },
  data() {
    return {
      eventTotal: 1000,
      allowedTotal: 1000,
      deniedTotal: 1000,
      unknownTotal: 1000,

      flag: false, // 判断是否显示图表组件
      categoryTotal: {}, // 每个分类下的文章数
      options: [],
      pieData: [],
      ip: {},
      user: {},
      opinfo: {},
      typeinfo: {}
    }
  },
  created() {
    this.getEventInfo(null, 0, 0)
    this.flag = true
  },
  methods: {
    fetchData() {
      console.log('Dashboard')
    },
    async getEventInfo(db, start, end) {
      const query = {}
      if (start === 0 && end === 0) {
        query.start = new Date().getTime() - 3600 * 1000 * 24 * 7
        query.end = new Date().getTime()
      } else {
        query.start = start
        query.end = end
      }
      query.db = db

      await api.getEventInfo(query).then((response) => {
        const dbs = response.data.db
        if (query.db == null || query.db.length === 0) {
          this.options = []
          for (const db of dbs) {
            const item = {}
            item.label = db
            item.value = db
            this.options.push(item)
          }
        }

        const types = response.data.types
        this.allowedTotal = types['1'] ? types['1'] : 0
        this.deniedTotal = types['2'] ? types['2'] : 0
        this.unknownTotal = types['3'] ? types['3'] : 0

        this.eventTotal = this.allowedTotal + this.deniedTotal + this.unknownTotal

        {
          const ops = response.data.ops
          this.pieData = []
          for (var i = 0; i <= 8; i++) {
            const index = i.toString()
            const item = {}
            item.name = sqlStrOpMap[index]
            item.value = ops[index] ? ops[index] : 0
            this.pieData.push(item)
          }
        }

        this.ip = response.data.ip
        this.user = response.data.user

        this.opinfo = response.data.opinfo
        this.opinfo.start = query.start
        this.opinfo.end = query.end

        this.typeinfo = response.data.typeinfo
        this.typeinfo.start = query.start
        this.typeinfo.end = query.end
      })
    },

    async queryData(name, start, end) {
      await this.getEventInfo(name, start, end)
    }
  }
}
</script>

<style lang="scss" scoped>
.dashboard {
  &-container {
    margin: 30px;
  }
  &-text {
    font-size: 30px;
    line-height: 46px;
  }
}
</style>
