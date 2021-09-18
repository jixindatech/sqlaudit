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
    },
    data: {
      type: Object,
      default: function() {
        return {}
      }
    }
  },
  data() {
    return {
      chart: null, // 引用echarts实例属性
      item: [],
      num: [],
      category: ['192.168.3.100', '192.168.123.100', '191.165.3.100', '10.0.0.100', '191.165.3.101', '191.165.3.102', '191.165.3.103', '191.165.3.10', '191.165.3.104', '191.165.3.105', '191.165.3.16', '191.165.3.107', '191.165.3.18'],
      barData: [3100, 2142, 1218, 581, 431, 383, 163, 2142, 1218, 581, 581, 431, 383]
    }
  },
  watch: {
    data: {
      handler(newValue, oldValue) {
        this.data = newValue
        this.item = this.data.item
        this.num = this.data.num
        this.initIPLineChart()
      },
      deep: true
    }
  },

  mounted() {
    this.$nextTick(() => {
      this.initIPLineChart()
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
    async initIPLineChart() {
      // await api.getEventsTopInfo().then((response) => {
      //  this.xAxisData = response.data.hosts
      //  this.seriesData = response.data.count
      // })
      this.chart = echarts.init(this.$refs.main, 'macarons')
      this.chart.setOption({
        title: { // 标题
          text: '请求IP汇总',
          left: 'center'
        },
        tooltip: {
          trigger: 'axis',
          axisPointer: {
            type: 'shadow'
          }
        },
        grid: {
          left: '3%',
          right: '4%',
          bottom: '3%',
          containLabel: true
        },
        xAxis: {
          type: 'value',
          axisLine: {
            show: false
          },
          axisTick: {
            show: false
          }
        },
        yAxis: {
          type: 'category',
          data: this.item,
          splitLine: { show: false },
          axisLine: {
            show: false
          },
          axisTick: {
            show: false
          },
          offset: 10,
          nameTextStyle: {
            fontSize: 15
          }
        },
        series: [
          {
            name: '数量',
            type: 'bar',
            data: this.num,
            barWidth: 14,
            barGap: 10,
            smooth: true,
            label: {
              normal: {
                show: true,
                position: 'right',
                offset: [5, -2],
                textStyle: {
                  color: '#F68300',
                  fontSize: 13
                }
              }
            },
            itemStyle: {
              emphasis: {
                barBorderRadius: 7
              },
              normal: {
                barBorderRadius: 7,
                color: new echarts.graphic.LinearGradient(
                  0, 0, 1, 0,
                  [
                    { offset: 0, color: '#3977E6' },
                    { offset: 1, color: '#37BBF8' }

                  ]
                )
              }
            }
          }
        ] })
    }
  }
}
</script>
