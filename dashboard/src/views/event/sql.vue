<template>
  <div class="app-container">
    <!-- 条件查询 -->
    <el-form :inline="true" :model="query" size="mini" :rules="rules">
      <el-form-item label="数据库名称:">
        <el-input v-model.trim="query.db" />
      </el-form-item>
      <el-form-item label="SQL关键字:">
        <el-input v-model.trim="query.sql" />
      </el-form-item>
      <el-form-item label="SQL指纹:">
        <el-input v-model="query.fingerprint" />
      </el-form-item>
      <el-form-item label="命中规则名称:">
        <el-input v-model.trim="query.name" />
      </el-form-item>
      <el-form-item label="规则类型:">
        <el-select v-model.number="query.type" clearable>
          <el-option
            v-for="item in sqlTypeOptions"
            :key="item.key"
            :label="item.value"
            :value="item.key"
          />
        </el-select>
      </el-form-item>
      <el-form-item label="请求IP地址:" prop="ip">
        <el-input v-model="query.ip" />
      </el-form-item>
      <el-form-item label="用户名称:">
        <el-input v-model.trim="query.user" />
      </el-form-item>
      <el-form-item label="请求类型:">
        <el-select v-model.number="query.op" clearable>
          <el-option
            v-for="item in sqlOptions"
            :key="item.key"
            :label="item.value"
            :value="item.key"
          />
        </el-select>
      </el-form-item>
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

      <el-form-item>
        <el-button icon="el-icon-search" type="primary" @click="queryData">查询</el-button>
        <el-button icon="el-icon-refresh" @click="reload">重置</el-button>
      </el-form-item>
    </el-form>
    <el-table
      :data="list"
      stripe
      border
      style="width: 100%"
    >
      <el-table-column align="center" prop="_source.user" label="数据库用户" />
      <el-table-column align="center" prop="_source.src" label="用户IP" />
      <el-table-column align="center" prop="_source.db" label="数据库" />
      <el-table-column align="center" prop="_source.sql" label="SQL命令" />
      <el-table-column align="center" prop="_source.fingerprint" label="指纹" />
      <el-table-column align="center" prop="_source.op" label="请求类型">
        <template slot-scope="scope">
          {{ sqlStrOpMap[scope.row._source.op] }}
        </template>
      </el-table-column>
      <el-table-column align="center" prop="_source.type" label="规则类型">
        <template slot-scope="scope">
          <el-tag v-show="scope.row._source.type===1" :type="scope.row._source.type | typeFilter">
            {{ '允许类型' }}
          </el-tag>
          <el-tag v-show="scope.row._source.type===2" :type="scope.row._source.type | typeFilter">
            {{ '拒绝类型' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column fixed align="center" prop="_source.time" label="时间">
        <template slot-scope="scope">
          {{ new Date(scope.row._source.time * 1000).toLocaleString() }}
        </template>
      </el-table-column>
    </el-table>
    <el-pagination
      :current-page="page.current"
      :page-sizes="[10, 20, 50]"
      :page-size="page.size"
      layout="total, sizes, prev, pager, next, jumper"
      :total="page.total"
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
    />
  </div>

</template>

<script>
import api from '@/api/event'
import { sqlTypeOptions, sqlStrOpMap, sqlOptions } from '@/utils/const'
import { validIP } from '@/utils/validate'
export default {
  filters: {
    typeFilter(status) {
      const statusMap = { 1: 'success', 2: 'danger' }
      return statusMap[status]
    }
  },
  data() {
    const validateIP = (rule, value, callback) => {
      if (value !== undefined && value !== '' && !validIP(value)) {
        callback(new Error('请输入正确的IP地址'))
      } else {
        callback()
      }
    }

    return {
      sqlOptions,
      sqlStrOpMap,
      sqlTypeOptions,
      list: [],
      listLoading: true,
      page: { // 分页对象
        current: 1, // 当前页码
        size: 20, // 每页显示多少条
        total: 0 // 总记录数
      },
      query: {},
      queryTime: [],
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
      rules: {
        ip: [
          { trigger: 'change', validator: validateIP }
        ]
      }
    }
  },
  created() {
    this.queryTime[0] = new Date().getTime() - 3600 * 1000 * 24 * 7
    this.queryTime[1] = new Date().getTime()
    this.fetchData()
  },
  methods: {
    async fetchData() {
      if (this.queryTime.length > 0) {
        this.query.start = this.queryTime[0]
        this.query.end = this.queryTime[1]
      }
      const { data } = await api.getList(this.query, this.page.current, this.page.size)
      this.list = data.data
      this.page.total = data.count
    },

    queryData() {
      this.page.current = 1
      this.fetchData()
    },
    reload() {
      this.query = {}
      this.queryTime[0] = new Date().getTime() - 3600 * 1000 * 24 * 7
      this.queryTime[1] = new Date().getTime()
      this.fetchData()
    },
    remoteClose() {
      this.edit.formData = {}
      this.edit.visible = false
      this.fetchData()
    },
    handleSizeChange(val) {
      this.page.size = val
      this.fetchData()
    },
    handleCurrentChange(val) {
      this.page.current = val
      this.fetchData()
    }
  }
}
</script>
