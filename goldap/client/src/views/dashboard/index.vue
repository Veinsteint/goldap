<template>
  <div class="dashboard-container">
    <div class="dashboard-editor-container">
      <!-- <github-corner class="github-corner" /> -->
      <panel-group @handleSetLineChartData="handleSetLineChartData" />
      <!-- <el-row style="background:#fff;padding:16px 16px 0;margin-bottom:32px;">
        <line-chart :chart-data="lineChartData" />
      </el-row> -->
    </div>
  </div>
</template>

<script>
// import GithubCorner from '@/components/GithubCorner'
import PanelGroup from './components/PanelGroup'
import { mapGetters } from 'vuex'
// import { Message } from 'element-ui'

export default {
  name: 'Dashboard',
  components: {
    // GithubCorner,
    PanelGroup,
    // eslint-disable-next-line vue/no-unused-components
  },
  computed: {
    ...mapGetters(['roles'])
  },
  data() {
    return {
      lineChartData: {}
    }
  },
  //普通用户登录后跳转到个人中心
 beforeRouteEnter(to, from, next) {
    next(vm => {
      const roles = vm.$store.getters.roles;
      if (roles.length > 0 && roles.includes('普通用户')) {
        vm.$router.push('/profile/index');
      }
    });
  },

   methods: {
    handleSetLineChartData(type) {
      this.lineChartData = {}
    }
  }
}
</script>

<style rel="stylesheet/scss" lang="scss" scoped>
  .dashboard-editor-container {
    padding: 32px;
    background-color: rgb(240, 242, 245);
    position: relative;

    .github-corner {
      position: absolute;
      top: 0;
      border: 0;
      right: 0;
    }

    .chart-wrapper {
      background: #fff;
      padding: 16px 16px 0;
      margin-bottom: 32px;
    }
  }

  @media (max-width:1024px) {
    .chart-wrapper {
      padding: 8px;
    }
  }
</style>
