<template>
  <c-modal-wrapper close>
    <template #title>
      <span>{{ title }}</span>
    </template>
    <template #text>
      <entity-dependencies-list-component
        :widget="widget"
        :columns="widget.parameters.widgetColumns"
        :entity-id="entity._id"
        :impact="config.impact"
      />
    </template>
  </c-modal-wrapper>
</template>

<script setup>
import { computed } from 'vue';

import { MODALS } from '@/constants';

import { generatePreparedDefaultContextWidget } from '@/helpers/entities/widget/form';

import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';

import EntityDependenciesListComponent from '@/components/other/entity/entity-dependencies-list.vue';

defineOptions({ name: MODALS.entityDependenciesList });

const props = defineProps({
  modal: {
    type: Object,
    required: true,
  },
});

const { t } = useI18n();
const { config } = useInnerModal(props);

const entity = computed(() => config.value.entity ?? {});
const title = computed(
  () =>
    config.value.title ??
    t('modals.entityDependenciesList.title', {
      name: entity.value.name,
    }),
);
const widget = computed(() => config.value.widget ?? generatePreparedDefaultContextWidget());
</script>
