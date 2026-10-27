<template>
  <c-modal-wrapper close>
    <template #title>
      {{ title }}
    </template>
    <template #text>
      <entity-upstream-network-graph :entity="entity" />
    </template>
  </c-modal-wrapper>
</template>

<script setup>
import { computed } from 'vue';

import { MODALS } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';

import EntityUpstreamNetworkGraph from '@/components/other/entity/entity-upstream-network-graph.vue';

defineOptions({ name: MODALS.entityUpstream });

const props = defineProps({
  modal: {
    type: Object,
    required: true,
  },
});

const { t } = useI18n();
const { config } = useInnerModal(props);

const entity = computed(() => config.value.entity ?? {});
const hasUpstream = computed(() => entity.value.upstream !== undefined);
const title = computed(
  () => config.value.title ?? (hasUpstream.value
    ? t('modals.entityUpstream.topLevelEntities')
    : t('modals.entityUpstream.entities')),
);
</script>
