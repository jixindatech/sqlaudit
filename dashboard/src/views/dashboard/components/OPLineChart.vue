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
      unknownData: [140, 122, 90, 280, 79, 78, 150, 45, 23, 120],
      selectData: [130, 140, 141, 142, 145, 150, 160, 11, 87, 20],
      unionData: [120, 82, 91, 154, 162, 140, 130, 140, 141, 142],
      insertData: [113, 72, 101, 34, 23, 12, 67, 12, 71, 127],
      updateData: [10, 12, 71, 127, 234, 45, 89, 43, 67, 120],
      deleteData: [130, 52, 21, 63, 45, 23, 120, 45, 34, 22],
      ddlData: [20, 62, 15, 124, 43, 67, 120, 62, 15, 124],
      showData: [190, 82, 20, 124, 56, 12, 230, 127, 234, 45],
      truncateData: [200, 32, 77, 112, 45, 34, 22, 154, 162, 140],
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
        this.unknownData = this.data['0'] ? this.data['0'] : defaultData
        this.selectData = this.data['1'] ? this.data['1'] : defaultData
        this.unionData = this.data['2'] ? this.data['2'] : defaultData
        this.insertData = this.data['3'] ? this.data['3'] : defaultData
        this.updateData = this.data['4'] ? this.data['4'] : defaultData
        this.deleteData = this.data['5'] ? this.data['5'] : defaultData
        this.ddlData = this.data['6'] ? this.data['6'] : defaultData
        this.showData = this.data['7'] ? this.data['7'] : defaultData
        this.truncateData = this.data['8'] ? this.data['8'] : defaultData

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
          text: 'SQL分布',
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
            interval: 0,
            rotate: -30
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
          data: ['SELECT', 'UNION', 'INSERT', 'UPDATE', 'DELETE', 'DDL', 'SHOW', 'TRUNCATE', 'UNKNOWN']
        },
        series: [
          {
            name: 'UNKNOWN', itemStyle: {
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
            data: this.unknownData,
            animationDuration: 2800,
            animationEasing: 'cubicInOut'
          },
          {
            name: 'SELECT',
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
            data: this.selectData,
            animationDuration: 2800,
            animationEasing: 'quadraticOut'
          },
          {
            name: 'UNION',
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
            data: this.unionData,
            animationDuration: 2800,
            animationEasing: 'quadraticOut'
          },
          {
            name: 'INSERT',
            smooth: true,
            type: 'line',
            itemStyle: {
              normal: {
                color: '#9000fa',
                lineStyle: {
                  color: '#9000fa',
                  width: 2
                },
                areaStyle: {
                  color: '#f3f8ff'
                }
              }
            },
            data: this.insertData,
            animationDuration: 2800,
            animationEasing: 'quadraticOut'
          },
          {
            name: 'UPDATE',
            smooth: true,
            type: 'line',
            itemStyle: {
              normal: {
                color: '#400000',
                lineStyle: {
                  color: '#400000',
                  width: 2
                },
                areaStyle: {
                  color: '#f3f8ff'
                }
              }
            },
            data: this.updateData,
            animationDuration: 2800,
            animationEasing: 'quadraticOut'
          },
          {
            name: 'DELETE',
            smooth: true,
            type: 'line',
            itemStyle: {
              normal: {
                color: '#666666',
                lineStyle: {
                  color: '#666666',
                  width: 2
                },
                areaStyle: {
                  color: '#f3f8ff'
                }
              }
            },
            data: this.deleteData,
            animationDuration: 2800,
            animationEasing: 'quadraticOut'
          },
          {
            name: 'DDL', itemStyle: {
              normal: {
                color: '#FF005A',
                lineStyle: {
                  color: '#FF005A',
                  width: 2
                }
              }
            },
            smooth: true,
            type: 'line',
            data: this.ddlData,
            animationDuration: 2800,
            animationEasing: 'cubicInOut'
          },
          {
            name: 'SHOW', itemStyle: {
              normal: {
                color: '#FFFF55',
                lineStyle: {
                  color: '#FFFF55',
                  width: 2
                }
              }
            },
            smooth: true,
            type: 'line',
            data: this.showData,
            animationDuration: 2800,
            animationEasing: 'cubicInOut'
          },
          {
            name: 'TRUNCATE', itemStyle: {
              normal: {
                color: '#00FF5A',
                lineStyle: {
                  color: '#00FF5A',
                  width: 2
                }
              }
            },
            smooth: true,
            type: 'line',
            data: this.truncateData,
            animationDuration: 2800,
            animationEasing: 'cubicInOut'
          }
        ]
      })
    }
  }
}
</script>
