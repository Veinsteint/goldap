<template>
  <section class="app-main" :style="appMainStyle">
    <transition name="fade-transform" mode="out-in">
      <keep-alive :include="cachedViews">
        <router-view :key="key" />
      </keep-alive>
    </transition>
    <el-footer class="footer-copyright" :style="footerStyle">
      <div>
        <span>Since 2025 </span>
        <el-divider direction="vertical" />
        <span>Powered by </span>
          <span>
            <a href="https://www.cimrbj.ac.cn/" target="_blank">CIMR</a>
          </span>
        <el-divider direction="vertical" />
        <span>Copyright </span>
          <span>
            <a href="https://github.com/cmplab-cimr" target="_blank">CMPLab</a>
          </span>
      </div>
    </el-footer>
  </section>
</template>

<script>
import { NAVBAR_HEIGHT, TAGSVIEW_HEIGHT, HEADER_HEIGHT_WITH_TAGS, FOOTER_COPYRIGHT_HEIGHT } from '@/utils/layout'

export default {
  name: 'AppMain',
  computed: {
    cachedViews() {
      return this.$store.state.tagsView.cachedViews
    },
    key() {
      return this.$route.path
    },
    footerStyle(){
      return {
        '--copyright-height': `${FOOTER_COPYRIGHT_HEIGHT}px`,
        height: `${FOOTER_COPYRIGHT_HEIGHT}px`
      }
    },
    appMainStyle() {
      return {
        '--navbar-height': `${NAVBAR_HEIGHT}px`,
        '--tagsview-height': `${TAGSVIEW_HEIGHT}px`,
        '--header-height-with-tags': `${HEADER_HEIGHT_WITH_TAGS}px`
      }
    }
  }
}
</script>

<style lang="scss" scoped>
.app-main {
  min-height: calc(100vh - var(--navbar-height, 50px));
  width: 100%;
  position: relative;
  overflow: hidden;
  padding-bottom: 1px; /* 底部版权信息 */
}

.fixed-header+.app-main {
  padding-top: 10px;
}

.footer-copyright {
  min-height: var(--copyright-height, 20px);
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  line-height: 20px;
  text-align: center;
  border-top: 1px dashed #dcdfe6;
  background: #fff;
  z-index: 1000;

  div {
    color: #999;
    font-size: 12px;
    font-weight: normal;
    
    a {
      color: #999;
      text-decoration: none;
      transition: color 0.3s;
      
      &:hover {
        color: #409eff;
      }
    }
  }
}

.hasTagsView {
  .app-main {
    min-height: calc(100vh - var(--header-height-with-tags, 84px));
    padding-bottom: 1px; /* 底部版权信息 */
  }

  .fixed-header+.app-main {
    padding-top: var(--header-height-with-tags, 84px);
  }
}
</style>

<style lang="scss">
// fix css style bug in open el-dialog
.el-popup-parent--hidden {
  .fixed-header {
    padding-right: 15px;
  }
}
</style>
