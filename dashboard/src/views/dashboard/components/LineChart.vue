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
    }
    /*
    chartData: {
      type: Object,
      required: true
    }
    */
  },
  data() {
    return {
      chart: null,
      selectData: [130, 140, 141, 142, 145, 150, 160, 11, 87, 20],
      unionData: [120, 82, 91, 154, 162, 140, 130, 140, 141, 142],
      insertData: [113, 72, 101, 34, 23, 12, 67, 12, 71, 127],
      updateData: [10, 12, 71, 127, 234, 45, 89, 43, 67, 120],
      deleteData: [130, 52, 21, 63, 45, 23, 120, 45, 34, 22],
      ddlData: [20, 62, 15, 124, 43, 67, 120, 62, 15, 124],
      showData: [190, 82, 20, 124, 56, 12, 230, 127, 234, 45],
      truncateData: [200, 32, 77, 112, 45, 34, 22, 154, 162, 140],
      unknownData: [140, 122, 90, 280, 79, 78, 150, 45, 23, 120]
    }
  },
  /*
  watch: {
    chartData: {
      deep: true,
      handler(val) {
        this.setOptions(val)
      }
    }
  },
  */
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
          data: ['1', '2', '3', '4', '5', '6', '7', '8', '9', '10'],
          boundaryGap: false,
          axisTick: {
            show: false
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
            data: this.actualData,
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
          },
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
          }
        ]
      })
    }
  }
}
</script>
