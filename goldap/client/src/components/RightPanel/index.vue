<template>
  <div ref="rightPanel" :class="{show:show}" class="rightPanel-container">
    <div class="rightPanel-background" />
    <div class="rightPanel">
      <div 
        class="handle-button-wrapper"
        :style="{'top':currentTop+'px'}"
        @mouseenter="isHovered = true"
        @mouseleave="isHovered = false"
      >
        <div 
          ref="handleButton"
          class="handle-button" 
          :class="{'is-hovered': isHovered}"
          :style="{'background-color':theme}" 
          @click="handleClick"
          @mousedown="handleMouseDown"
        >
          <i :class="show?'el-icon-close':'el-icon-setting'" />
        </div>
      </div>
      <div class="rightPanel-items">
        <slot />
      </div>
    </div>
  </div>
</template>

<script>
import { addClass, removeClass } from '@/utils'
import { getRightPanelTop, RIGHT_PANEL_DEFAULT_TOP } from '@/utils/layout'

export default {
  name: 'RightPanel',
  props: {
    clickNotClose: {
      default: false,
      type: Boolean
    },
    buttonTop: {
      default: RIGHT_PANEL_DEFAULT_TOP,
      type: Number
    }
  },
  data() {
    return {
      show: false,
      isHovered: false,
      currentTop: this.buttonTop,
      isDragging: false,
      dragStartY: 0,
      dragStartTop: 0,
      hasMoved: false
    }
  },
  computed: {
    theme() {
      return this.$store.state.settings.theme
    }
  },
  watch: {
    show(value) {
      if (value && !this.clickNotClose) {
        this.addEventClick()
      }
      if (value) {
        addClass(document.body, 'showRightPanel')
      } else {
        removeClass(document.body, 'showRightPanel')
      }
    },
    buttonTop(newVal) {
      if (!this.isDragging) {
        this.currentTop = newVal
      }
    }
  },
  mounted() {
    this.insertToBody()
    this.currentTop = this.buttonTop
    window.addEventListener('mousemove', this.handleMouseMove)
    window.addEventListener('mouseup', this.handleMouseUp)
  },
  beforeDestroy() {
    const elx = this.$refs.rightPanel
    elx.remove()
    window.removeEventListener('mousemove', this.handleMouseMove)
    window.removeEventListener('mouseup', this.handleMouseUp)
  },
  methods: {
    addEventClick() {
      window.addEventListener('click', this.closeSidebar)
    },
    closeSidebar(evt) {
      const parent = evt.target.closest('.rightPanel')
      if (!parent) {
        this.show = false
        window.removeEventListener('click', this.closeSidebar)
      }
    },
    insertToBody() {
      const elx = this.$refs.rightPanel
      const body = document.querySelector('body')
      body.insertBefore(elx, body.firstChild)
    },
    handleClick(e) {
      // 如果发生了拖动，则不触发点击事件
      if (this.hasMoved) {
        e.preventDefault()
        e.stopPropagation()
        return
      }
      this.show = !this.show
    },
    handleMouseDown(e) {
      e.preventDefault()
      e.stopPropagation()
      this.isDragging = true
      this.hasMoved = false
      this.dragStartY = e.clientY
      this.dragStartTop = this.currentTop
      document.body.style.userSelect = 'none'
      document.body.style.cursor = 'move'
    },
    handleMouseMove(e) {
      if (!this.isDragging) return
      
      const deltaY = e.clientY - this.dragStartY
      
      // 如果移动距离超过5px，认为是拖动操作
      if (Math.abs(deltaY) > 5) {
        this.hasMoved = true
      }
      
      const newTop = this.dragStartTop + deltaY
      
      // 限制拖动范围，确保按钮在可视区域内
      const minTop = 0
      const maxTop = window.innerHeight - 48 // 48px 是按钮高度
      
      this.currentTop = Math.max(minTop, Math.min(maxTop, newTop))
      
      // 触发事件通知父组件
      this.$emit('top-change', this.currentTop)
    },
    handleMouseUp() {
      if (this.isDragging) {
        this.isDragging = false
        document.body.style.userSelect = ''
        document.body.style.cursor = ''
        // 延迟重置 hasMoved，确保 click 事件能正确判断
        setTimeout(() => {
          this.hasMoved = false
        }, 0)
      }
    }
  }
}
</script>

<style>
.showRightPanel {
  overflow: hidden;
  position: relative;
  width: calc(100% - 15px);
}
</style>

<style lang="scss" scoped>
.rightPanel-background {
  position: fixed;
  top: 0;
  left: 0;
  opacity: 0;
  transition: opacity .3s cubic-bezier(.7, .3, .1, 1);
  background: rgba(0, 0, 0, .2);
  z-index: -1;
}

.rightPanel {
  width: 100%;
  max-width: 260px;
  height: 100vh;
  position: fixed;
  top: 0;
  right: 0;
  box-shadow: 0px 0px 15px 0px rgba(0, 0, 0, .05);
  transition: all .25s cubic-bezier(.7, .3, .1, 1);
  transform: translate(100%);
  background: #fff;
  z-index: 40000;
}

.show {
  transition: all .3s cubic-bezier(.7, .3, .1, 1);

  .rightPanel-background {
    z-index: 20000;
    opacity: 1;
    width: 100%;
    height: 100%;
  }

  .rightPanel {
    transform: translate(0);
  }
  
  .handle-button {
    left: -48px !important;
  }
}

.handle-button-wrapper {
  position: absolute;
  left: -18;
  width: 58px; 
  height: 48px;
  z-index: 0;
  pointer-events: auto;
}

.handle-button {
  width: 48px;
  height: 48px;
  position: absolute;
  left: -8px;
  text-align: center;
  font-size: 24px;
  border-radius: 6px 0 0 6px !important;
  z-index: 0;
  pointer-events: auto;
  cursor: pointer;
  color: #fff;
  line-height: 48px;
  transition: left 0.2s cubic-bezier(.7, .3, .1, 1);
  user-select: none;
  
  i {
    font-size: 24px;
    line-height: 48px;
  }
  
  &.is-hovered {
    left: -48px;
  }
}
</style>
