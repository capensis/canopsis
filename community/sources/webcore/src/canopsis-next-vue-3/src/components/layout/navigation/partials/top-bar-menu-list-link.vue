<template>
  <v-list-item
    :class="link.class"
    :disabled="disabled"
    class="top-bar-menu-link"
    @click="click"
    @mouseenter="handleMouseEnter"
  >
    <template #prepend>
      <v-icon
        :color="link.color || undefined"
        size="24"
      >
        {{ link.icon }}
      </v-icon>
    </template>

    <v-list-item-title :class="{ [`text-${link.color}`]: link.color }">
      <c-layout
        align="center"
        justify="space-between"
      >
        <span>{{ link.title }}</span>
      </c-layout>
    </v-list-item-title>

    <template
      v-if="link.links"
      #append
    >
      <v-icon>arrow_right</v-icon>
    </template>
  </v-list-item>
</template>

<script>
import { computed } from 'vue'
import { useRouter } from 'vue-router'

export default {
  props: {
    link: {
      type: Object,
      default: () => ({}),
    },
  },

  emits: ['click', 'mouseenter'],

  setup(props, { emit }) {
    const router = useRouter()

    const disabled = computed(() => (
      props.link.route?.name && !router.hasRoute(props.link.route.name)
    ))

    const click = async () => {
      if (props.link.links?.length) {
        return
      }

      if (props.link.handler) {
        await props.link.handler()
        emit('click')

        return
      }

      if (props.link.route && !disabled.value) {
        await router.push(props.link.route)
        emit('click')
      }
    }

    const handleMouseEnter = event => emit('mouseenter', event)

    return {
      disabled,
      click,
      handleMouseEnter,
    }
  },
}
</script>

<style lang="scss" scoped>
.top-bar-menu-link {
  min-width: 190px;
  text-decoration: none;
}
</style>
