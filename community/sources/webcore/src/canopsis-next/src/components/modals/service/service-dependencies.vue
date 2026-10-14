<template>
  <c-modal-wrapper close>
    <template #title>
      <span>{{ $t('modals.serviceDependencies.title') }}</span>
    </template>
    <template #text>
      <service-dependencies
        :root="config.root"
        :columns="config.columns"
        :impact="config.impact"
        include-root
      />
    </template>
  </c-modal-wrapper>
</template>

<script>
import ServiceDependencies from '@/components/other/service/partials/service-dependencies.vue';

export default {
  components: { ServiceDependencies },

  props: {
    modal: {
      type: Object,
      required: true,
    },
  },

  computed: {
    config() {
      return this.modal.config ?? {};
    },
  },
};
</script>
