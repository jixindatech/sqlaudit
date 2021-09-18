<template>
  <div ref="main" :class="className" :style="{height: height, width: width}" />
</template>

<script>

// import echarts from 'echarts'
import * as echarts from 'echarts'
require('echarts/theme/macarons')

import resize from './mixins/resize'
// import api from '@/api/event'

export default {
  mixins: [resize],
  props: {
    className: {
      type: String,
      default: 'chart'
    },
    width: {
      type: String,
      default: '100%'
    },
    height: {
      type: String,
      default: '400px'
    }
  },
  data() {
    return {
      chart: null, // 引用echarts实例属性
      xAxisData: [4, 5, 6],
      seriesData: [4, 5, 6]
    }
  },
  mounted() {
    this.$nextTick(() => {
      this.initBarChart()
    })
  },
  beforeDestroy() {
    if (!this.chart) {
      return
    }
    this.chart.dispose()
    this.chart = null
  },
  methods: {
    async initBarChart() {
      // await api.getEventsTopInfo().then((response) => {
      //  this.xAxisData = response.data.hosts
      //  this.seriesData = response.data.count
      // })
      this.chart = echarts.init(this.$refs.main, 'macarons')
      this.chart.setOption({
        // color: ['#3398DB'],
        title: { // 标题
          text: '客户端请求汇总',
          left: 'center'
        },
        tooltip: { // 提示框
          trigger: 'axis',
          axisPointer: { // 坐标轴指示器，坐标轴触发有效
            type: 'shadow' // 默认为直线，可选为：'line' | 'shadow'
          }
        },
        grid: { // 柱状图整体位置
          left: '3%',
          right: '4%',
          bottom: '3%',
          containLabel: true
        },
        xAxis: [ // x轴
          {
            type: 'category',
            data: this.xAxisData,
            axisTick: {
              alignWithLabel: true // 类目轴中在 boundaryGap 为 true 的时候有效，可以保证刻度线和标签对齐
            }
          }
        ],
        yAxis: [ // y 轴
          {
            type: 'value'
          }
        ],
        series: [
          {
            name: '事件数',
            type: 'bar', // 柱状图
            barWidth: '20%', // 每个柱子宽度
            data: this.seriesData
          }
        ]
      })
    }
  }
}
</script>
