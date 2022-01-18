<template>
  <div>
    <el-form :inline="true" :model="query" size="mini">
      <el-form-item>
        <el-date-picker
          v-model="queryTime"
          type="datetimerange"
          :picker-options="pickerOptions"
          range-separator="-"
          start-placeholder=""
          end-placeholder=""
          value-format="timestamp"
          align="right"
        />
      </el-form-item>
      <el-form-item label="数据库名称:">
        <el-select v-model="query.db" clearable placeholder="全部数据库">
          <el-option
            v-for="item in options"
            :key="item.value"
            :label="item.label"
            :value="item.value"
          />
        </el-select>
      </el-form-item>
      <el-form-item>
        <el-button icon="el-icon-search" type="primary" @click="queryFunction(query, queryTime)">查询</el-button>
        <el-button icon="el-icon-refresh" @click="reloadFunction">重置</el-button>
      </el-form-item>
    </el-form>

    <el-row :gutter="40" class="panel-group">
      <el-col :xs="12" :sm="12" :lg="5" class="card-panel-col">
        <div class="card-panel">
          <div class="card-panel-icon-wrapper icon-people">
            <svg-icon icon-class="tree-table" class-name="card-panel-icon" />
          </div>
          <div class="card-panel-description">
            <div class="card-panel-text">
              事件总数
            </div>
            <count-to :start-val="0" :end-val="eventTotal" :duration="2600" class="card-panel-num" />
          </div>
        </div>
      </el-col>
      <el-col :xs="12" :sm="12" :lg="5" class="card-panel-col">
        <div class="card-panel">
          <div class="card-panel-icon-wrapper icon-people">
            <svg-icon icon-class="tree-table" class-name="card-panel-icon" />
          </div>
          <div class="card-panel-description">
            <div class="card-panel-text">
              指纹类型
            </div>
            <count-to :start-val="0" :end-val="fingerprintTotal" :duration="2600" class="card-panel-num" />
          </div>
        </div>
      </el-col>
      <el-col :xs="12" :sm="12" :lg="4" class="card-panel-col">
        <div class="card-panel">
          <div class="card-panel-icon-wrapper icon-people">
            <svg-icon icon-class="tree-table" class-name="card-panel-icon" />
          </div>
          <div class="card-panel-description">
            <div class="card-panel-text">
              拒绝类型
            </div>
            <count-to :start-val="0" :end-val="deniedTotal" :duration="2600" class="card-panel-num" />
          </div>
        </div>
      </el-col>

      <el-col :xs="12" :sm="12" :lg="4" class="card-panel-col">
        <div class="card-panel">
          <div class="card-panel-icon-wrapper icon-message">
            <svg-icon icon-class="documentation" class-name="card-panel-icon" />
          </div>
          <div class="card-panel-description">
            <div class="card-panel-text">
              允许类型
            </div>
            <count-to :start-val="0" :end-val="allowedTotal" :duration="3000" class="card-panel-num" />
          </div>
        </div>
      </el-col>
      <el-col :xs="12" :sm="12" :lg="4" class="card-panel-col">
        <div class="card-panel">
          <div class="card-panel-icon-wrapper icon-money">
            <svg-icon icon-class="question" class-name="card-panel-icon" />
          </div>
          <div class="card-panel-description">
            <div class="card-panel-text">
              未知类型
            </div>
            <count-to :start-val="0" :end-val="unknownTotal" :duration="3200" class="card-panel-num" />
          </div>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import CountTo from 'vue-count-to'

export default {
  components: {
    CountTo
  },

  props: {
    eventTotal: {
      type: Number,
      default: 0
    },
    fingerprintTotal: {
      type: Number,
      default: 0
    },
    allowedTotal: {
      type: Number,
      default: 0
    },
    deniedTotal: {
      type: Number,
      default: 0
    },
    unknownTotal: {
      type: Number,
      default: 0
    },

    options: {
      type: Array,
      default: () => []
    },
    queryData: {
      type: Function,
      default: () => {}
    },
    reload: {
      type: Function,
      default: () => {}
    }
  },
  data() {
    return {
      query: {
      },
      pickerOptions: {
        shortcuts: [{
          text: '最近30分钟',
          onClick(picker) {
            const end = new Date()
            const start = new Date()
            start.setTime(start.getTime() - 1800 * 1000)
            picker.$emit('pick', [start, end])
          }
        }, {
          text: '最近一小时',
          onClick(picker) {
            const end = new Date()
            const start = new Date()
            start.setTime(start.getTime() - 3600 * 1000)
            picker.$emit('pick', [start, end])
          }
        }, {
          text: '最近24小时',
          onClick(picker) {
            const end = new Date()
            const start = new Date()
            start.setTime(start.getTime() - 3600 * 1000 * 24)
            picker.$emit('pick', [start, end])
          }
        },
        {
          text: '最近一周',
          onClick(picker) {
            const end = new Date()
            const start = new Date()
            start.setTime(start.getTime() - 3600 * 1000 * 24 * 7)
            picker.$emit('pick', [start, end])
          }
        },
        {
          text: '最近一个月',
          onClick(picker) {
            const end = new Date()
            const start = new Date()
            start.setTime(start.getTime() - 3600 * 1000 * 24 * 30)
            picker.$emit('pick', [start, end])
          }
        }]
      },
      queryTime: []
    }
  },
  created() {
    this.initDefaultQueryTime()
  },
  methods: {
    initDefaultQueryTime() {
      this.queryTime[0] = new Date().getTime() - 3600 * 1000 * 24 * 7
      this.queryTime[1] = new Date().getTime()
    },
    queryFunction(query, timeQuery) {
      this.queryData(query.db, timeQuery[0], timeQuery[1])
    },
    reloadFunction() {
      this.initDefaultQueryTime()
      this.query = {}
      this.queryData(null, this.queryTime[0], this.queryTime[1])
    }
  }
}
</script>

<style lang="scss" scoped>
.panel-group {
  margin-top: 18px;

  .card-panel-col {
    margin-bottom: 32px;
  }

  .card-panel {
    height: 108px;
    cursor: pointer;
    font-size: 12px;
    position: relative;
    overflow: hidden;
    color: #666;
    background: #fff;
    box-shadow: 4px 4px 40px rgba(0, 0, 0, .05);
    border-color: rgba(0, 0, 0, .05);

    &:hover {
      .card-panel-icon-wrapper {
        color: #fff;
      }

      .icon-people {
        background: #40c9c6;
      }

      .icon-message {
        background: #36a3f7;
      }

      .icon-money {
        background: #f4516c;
      }

      .icon-shopping {
        background: #34bfa3
      }
    }

    .icon-people {
      color: #40c9c6;
    }

    .icon-message {
      color: #36a3f7;
    }

    .icon-money {
      color: #f4516c;
    }

    .icon-shopping {
      color: #34bfa3
    }

    .card-panel-icon-wrapper {
      float: left;
      margin: 14px 0 0 14px;
      padding: 16px;
      transition: all 0.38s ease-out;
      border-radius: 6px;
    }

    .card-panel-icon {
      float: left;
      font-size: 48px;
    }

    .card-panel-description {
      float: right;
      font-weight: bold;
      margin-top: 26px;
      margin-bottom: 26px;
      margin-right: 10px;
      margin-left: 0px;

      .card-panel-text {
        line-height: 18px;
        color: rgba(0, 0, 0, 0.45);
        font-size: 16px;
        margin-bottom: 12px;
      }

      .card-panel-num {
        font-size: 20px;
      }
    }
  }
}

@media (max-width:550px) {
  .card-panel-description {
    display: none;
  }

  .card-panel-icon-wrapper {
    float: none !important;
    width: 100%;
    height: 100%;
    margin: 0 !important;

    .svg-icon {
      display: block;
      margin: 14px auto !important;
      float: none !important;
    }
  }
}
</style>
