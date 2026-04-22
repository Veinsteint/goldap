/**
 * 布局相关常量
 * 用于计算页面高度
 */

/**
 * 底部版权信息高度
 */
export const FOOTER_COPYRIGHT_HEIGHT = 20

// Navbar 高度
export const NAVBAR_HEIGHT = 50

// TagsView 高度
export const TAGSVIEW_HEIGHT = 34

// 计算总高度（当 tags-view 显示时）
export const HEADER_HEIGHT_WITH_TAGS = NAVBAR_HEIGHT + TAGSVIEW_HEIGHT // 84px

// 计算总高度（当 tags-view 隐藏时）
export const HEADER_HEIGHT_WITHOUT_TAGS = NAVBAR_HEIGHT // 50px

// RightPanel 默认位置（当 tags-view 显示时）
export const RIGHT_PANEL_DEFAULT_TOP = HEADER_HEIGHT_WITH_TAGS + 2

/**
 * 获取当前 header 高度
 * @param {boolean} hasTagsView - 是否显示 tags-view
 * @returns {number} header 总高度
 */
export function getHeaderHeight(hasTagsView = true) {
  return hasTagsView ? HEADER_HEIGHT_WITH_TAGS : HEADER_HEIGHT_WITHOUT_TAGS
}

/**
 * 获取右侧面板高度
 * @param {boolean} hasTagsView - 是否显示 tags-view
 * @returns {number} 右侧面板位置
 */
export function getRightPanelTop(hasTagsView = true) {
  return hasTagsView ? RIGHT_PANEL_DEFAULT_TOP : RIGHT_PANEL_DEFAULT_TOP - TAGSVIEW_HEIGHT
}

/**
 * 计算内容区域高度
 * @param {boolean} hasTagsView - 是否显示 tags-view
 * @param {number} padding - 额外的 padding 值（默认 0）
 * @returns {string} CSS calc 表达式
 */
export function getContentHeight(hasTagsView = true, padding = 0) {
  const headerHeight = getHeaderHeight(hasTagsView)
  if (padding > 0) {
    return `calc(100vh - ${headerHeight + padding}px)`
  }
  return `calc(100vh - ${headerHeight}px)`
}
