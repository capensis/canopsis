<script lang="ts">
import { h } from 'vue'

export default {
  inheritAttrs: false,

  props: {
    template: {
      type: String,
      required: false,
    },
    parent: {
      type: Object,
      required: false,
    },
    templateProps: {
      type: Object,
      default: () => ({}),
    },
  },

  computed: {
    parentNode() {
      return this.parent || this.$parent
    },

    parentOptions() {
      return this.parentNode?.$options ?? {}
    },

    methodProps() {
      return Object.entries(this.parentOptions.methods ?? {}).reduce((acc, [key]) => {
        acc[key] = this.parentNode[key]

        return acc
      }, {})
    },

    props() {
      return Object.entries({
        ...this.parentNode?.$data,
        ...this.parentNode?.$props,
        ...this.methodProps,
        ...this.templateProps,
      }).reduce((acc, [key, value]) => {
        if (!this.isDeclaredByCurrentComponent(key)) {
          acc[key] = value
        }

        return acc
      }, {})
    },

    propsTypes() {
      return Object.keys({
        ...this.parentNode?.$data,
        ...this.parentOptions.props,
        ...this.parentOptions.methods,
        ...this.templateProps,
      }).filter(key => !this.isDeclaredByCurrentComponent(key))
    },

    component() {
      return {
        template: this.template || '<div />',
        components: this.parentOptions.components,
        computed: this.parentOptions.computed,
        props: this.propsTypes,
        compilerOptions: {
          whitespace: 'preserve' as const,
        },

        setup: () => ({
          ...this.parentNode?.$.provides,
        }),
      }
    },
  },

  methods: {
    isDeclaredByCurrentComponent(key) {
      return (
        key in (this.$data ?? {})
        || key in (this.$props ?? {})
        || key in (this.$options.computed ?? {})
        || key in (this.$options.methods ?? {})
      )
    },
  },

  render() {
    if (!this.template) {
      return null
    }

    return h(this.component, {
      ...this.props,
      ...this.$attrs,
    })
  },
}
</script>
