<template>
  <div ref="main" :class="className" :style="{height:height,width:width}" />
</template>

<script>
import * as echarts from 'echarts'
require('echarts/theme/macarons') // echarts theme
import resize from './mixins/resize'

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
      default: '450px'
    },
    autoResize: {
      type: Boolean,
      default: true
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
      chart: null,
      allwoedData: [140, 122, 90, 280, 79, 78, 150, 45, 23, 120],
      deniedData: [130, 140, 141, 142, 145, 150, 160, 11, 87, 20],
      unknownData: [120, 82, 91, 154, 162, 140, 130, 140, 141, 142],
      timeData: []
    }
  },
  watch: {
    data: {
      handler(newValue, oldValue) {
        const defaultData = [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0]
        this.data = newValue
        this.item = this.data.item
        this.num = this.data.num
        this.allwoedData = this.data['1'] ? this.data['1'] : defaultData
        this.deniedData = this.data['2'] ? this.data['2'] : defaultData
        this.unknownData = this.data['3'] ? this.data['3'] : defaultData

        const intervalTime = (this.data.end - this.data.start) / 10
        const timesSplice = []
        for (var i = 0; i < 10; i++) {
          const timeData = this.data.start + i * intervalTime
          timesSplice.push((new Date(timeData)).toLocaleString())
        }
        timesSplice.push((new Date(this.data.end)).toLocaleString())
        this.timeData = timesSplice
        this.initChart()
      },
      deep: true
    }
  },
  mounted() {
    this.$nextTick(() => {
      this.initChart()
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
    initChart() {
      this.chart = echarts.init(this.$refs.main, 'macarons')
      this.setOptions(this.chartData)
    },
    setOptions({ expectedData, actualData } = {}) {
      this.chart.setOption({
        title: { // 标题
          text: '类型分布',
          left: 'left'
        },

        xAxis: {
          name: '日期',
          data: this.timeData,
          boundaryGap: false,
          axisTick: {
            show: false
          },
          axisLabel: {
            interval: 0
          }
        },
        grid: {
          left: 10,
          right: 10,
          bottom: 20,
          top: 30,
          containLabel: true
        },
        tooltip: {
          trigger: 'axis',
          axisPointer: {
            type: 'cross'
          },
          padding: [5, 10]
        },
        yAxis: {
          axisTick: {
            show: false
          }
        },
        legend: {
          data: ['Allowed', 'Denied', 'Unknown']
        },
        series: [
          {
            name: 'Allowed', itemStyle: {
              normal: {
                color: '#eeee5A',
                lineStyle: {
                  color: '#eeee5A',
                  width: 2
                }
              }
            },
            smooth: true,
            type: 'line',
            data: this.allwoedData,
            animationDuration: 2800,
            animationEasing: 'cubicInOut'
          },
          {
            name: 'Denied',
            smooth: true,
            type: 'line',
            itemStyle: {
              normal: {
                color: '#3888fa',
                lineStyle: {
                  color: '#3888fa',
                  width: 2
                },
                areaStyle: {
                  color: '#f3f8ff'
                }
              }
            },
            data: this.deniedData,
            animationDuration: 2800,
            animationEasing: 'quadraticOut'
          },
          {
            name: 'Unknown',
            smooth: true,
            type: 'line',
            itemStyle: {
              normal: {
                color: '#2000ab',
                lineStyle: {
                  color: '#2000ab',
                  width: 2
                },
                areaStyle: {
                  color: '#f3f8ff'
                }
              }
            },
            data: this.unknownData,
            animationDuration: 2800,
            animationEasing: 'quadraticOut'
          }
        ]
      })
    }
  }
}
</script>
