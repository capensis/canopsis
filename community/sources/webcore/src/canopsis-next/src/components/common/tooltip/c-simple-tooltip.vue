<script>
import { isNumber } from 'lodash';
import { computed, watch, onBeforeUnmount } from 'vue';

import { getTooltipDimensions } from '@/helpers/tooltip/tooltip';

const TOOLTIP_CLASSES = {
  base: 'c-simple-tooltip__content',
  active: 'c-simple-tooltip__content--active',
  clickable: 'c-simple-tooltip__content--clickable',
  top: 'c-simple-tooltip__content--top',
  right: 'c-simple-tooltip__content--right',
  bottom: 'c-simple-tooltip__content--bottom',
  left: 'c-simple-tooltip__content--left',
};

export default {
  props: {
    content: {
      type: String,
      default: '',
    },
    top: {
      type: Boolean,
      default: false,
    },
    right: {
      type: Boolean,
      default: false,
    },
    bottom: {
      type: Boolean,
      default() {
        return !this.top && !this.right && !this.left;
      },
    },
    left: {
      type: Boolean,
      default: false,
    },
    nudgeTop: {
      type: Number,
      default: 10,
    },
    nudgeRight: {
      type: Number,
      default: 10,
    },
    nudgeBottom: {
      type: Number,
      default: 10,
    },
    nudgeLeft: {
      type: Number,
      default: 10,
    },
    contentMargin: {
      type: Number,
      default: 12,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    clickable: {
      type: Boolean,
      default: false,
    },
    maxWidth: {
      type: Number,
      required: false,
    },
  },
  setup(props) {
    let dimensions = {};
    let isActive = false;
    let activatorElement = null;
    let contentElement = null;

    /**
     * TOOLTIP CONTENT ELEMENT METHODS
     */
    const onTransitionEnd = () => {
      if (isActive) {
        return;
      }

      contentElement.style.display = 'none';
    };

    const setTooltipContent = (content) => {
      if (!contentElement) {
        return;
      }

      contentElement.innerHTML = content;
    };

    const hideTooltipContent = () => {
      if (!contentElement) {
        return;
      }

      activatorElement = null;
      isActive = false;
      contentElement.style.transform = 'translate(0, 0)';
      contentElement.classList.remove(TOOLTIP_CLASSES.active);
    };

    const onTooltipContentMouseLeave = (event) => {
      if (props.clickable && activatorElement?.contains(event.toElement)) {
        return;
      }

      hideTooltipContent();
    };

    const mountTooltipContent = () => {
      if (contentElement) {
        return;
      }

      contentElement = document.createElement('div');
      contentElement.classList.add(TOOLTIP_CLASSES.base);

      contentElement.addEventListener('transitionend', onTransitionEnd);

      if (props.clickable) {
        contentElement.addEventListener('mouseleave', onTooltipContentMouseLeave);

        const [className, nudge] = {
          [props.top]: [TOOLTIP_CLASSES.top, props.nudgeTop],
          [props.right]: [TOOLTIP_CLASSES.right, props.nudgeRight],
          [props.bottom]: [TOOLTIP_CLASSES.top, props.nudgeBottom],
          [props.left]: [TOOLTIP_CLASSES.left, props.nudgeLeft],
        }.true ?? [TOOLTIP_CLASSES.bottom, props.nudgeBottom];

        contentElement.classList.add(TOOLTIP_CLASSES.clickable, className);

        contentElement.style.setProperty('--tooltip-content-nudge', `${nudge}px`);
      }

      if (props.maxWidth) {
        contentElement.style.maxWidth = `${props.maxWidth}${isNumber(props.maxWidth) ? 'px' : ''}`;
      }

      contentElement.innerHTML = props.content;

      document.querySelector('[data-app]').appendChild(contentElement);
    };

    const showTooltipContent = () => {
      if (!contentElement) {
        mountTooltipContent();
      }

      isActive = true;
      contentElement.style.display = 'inline-block';
      contentElement.classList.add(TOOLTIP_CLASSES.active);
    };

    /**
     * ACTIVATOR LISTENERS
     */
    const onMouseEnter = (event) => {
      activatorElement = event.target;

      showTooltipContent();

      dimensions = getTooltipDimensions({
        ...props,

        contentElement,
        targetElement: event.target,
      });

      contentElement.style.top = `${dimensions.top}px`;
      contentElement.style.left = `${dimensions.left}px`;
      contentElement.style.zIndex = dimensions.zIndex;

      window.requestAnimationFrame(() => contentElement.style.transform = dimensions.transform);
    };

    const onMouseLeave = (event) => {
      if (props.clickable && contentElement?.contains(event.toElement)) {
        return;
      }

      hideTooltipContent();
    };

    const activatorListeners = computed(() => (
      props.disabled
        ? {}
        : {
          mouseenter: onMouseEnter,
          mouseleave: onMouseLeave,
        }
    ));

    /**
     * WATCHERS
     */
    watch(() => props.content, setTooltipContent);

    /**
     * LIFECYCLE HOOKS
     */
    onBeforeUnmount(() => contentElement?.remove());

    return { activatorListeners };
  },
  render() {
    return this.$scopedSlots.activator({ on: this.activatorListeners });
  },
};
</script>

<style lang="scss">
.c-simple-tooltip__content {
  --tooltip-content-nudge: 10px;

  position: absolute;
  display: inline-block;
  max-width: 60%;
  width: auto;
  opacity: 0;
  background: rgba(97, 97, 97, 0.9);
  color: #FFFFFF;
  border-radius: 2px;
  font-size: 12px;
  line-height: 19px;
  padding: 5px 8px;
  text-transform: initial;
  pointer-events: none;
  transition: opacity .3s cubic-bezier(0, 0, 0.2, 1), transform .3s cubic-bezier(0, 0, 0.2, 1);

  &--active {
    transition: transform .3s cubic-bezier(0, 0, 0.2, 1);
    opacity: .9;
  }

  &--clickable {
    pointer-events: initial;

    & > * {
      position: relative;
      z-index: 2;
    }

    &:after {
      left: 0;
      content: '';
      display: block;
      position: absolute;
      width: 100%;
      height: 100%;
      z-index: 1;
    }

    &.c-simple-tooltip__content--top:after {
      top: var(--tooltip-content-nudge);
    }

    &.c-simple-tooltip__content--right:after {
      right: var(--tooltip-content-nudge);
    }

    &.c-simple-tooltip__content--bottom:after {
      bottom: var(--tooltip-content-nudge);
    }

    &.c-simple-tooltip__content--left:after {
      left: var(--tooltip-content-nudge);
    }
  }
}
</style>
