<template>
  <v-tooltip location="left">
    <template #activator="{ props }">
      <v-btn
        v-bind="props"
        :active="activeViewEditing"
        :loading="activeViewEditingProcess"
        class="view-editing-btn__button"
        icon
        @click.stop="toggleEditing"
      >
        <v-icon
          :icon="activeViewEditing ? 'done' : 'edit'"
          size="small"
        />
      </v-btn>
    </template>

    <div>
      <div>{{ $t('common.toggleEditView') }} (ctrl + e / command + e)</div>
      <div class="font-italic">
        {{ $t('common.toggleEditViewSubtitle') }}
      </div>
    </div>
  </v-tooltip>
</template>

<script>
import { activeViewMixin } from '@/mixins/active-view';

export default {
  mixins: [activeViewMixin],
};
</script>

<style lang="scss" scoped>
.view-editing-btn__button {
  border-color: #212121 !important;
  background-color: #212121 !important;
  color: #fff !important;

  .v-theme--dark & {
    border-color: #979797 !important;
    background-color: #979797 !important;
  }
}
</style>
