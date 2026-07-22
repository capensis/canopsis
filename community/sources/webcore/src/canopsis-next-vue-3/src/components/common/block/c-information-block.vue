<template>
  <c-layout class="c-information-block" column>
    <c-layout v-if="title || helpText">
      <div class="text-subtitle-1 font-weight-bold">
        {{ title }}
      </div>
      <c-help-icon
        v-if="helpText"
        :text="helpText"
        :icon="helpIcon"
        :color="helpIconColor"
        icon-class="ml-2"
        max-width="300"
        top
      />
    </c-layout>
    <c-layout v-if="$slots.subtitle" class="mt-2">
      <span class="text-medium-emphasis">
        <slot name="subtitle" />
      </span>
    </c-layout>
    <slot />
  </c-layout>
</template>

<script>
// TODO: remove this component after legacy views are fully migrated.
export default {
  props: {
    title: {
      type: String,
      required: true,
    },
    helpText: {
      type: String,
      required: false,
    },
    helpIcon: {
      type: String,
      required: false,
    },
    helpIconColor: {
      type: String,
      required: false,
    },
  },
}
</script>
