<template>
  <div
    :class="[themeClasses]"
    class="c-advanced-search v-input v-input--hide-details theme--light
    v-text-field v-text-field--single-line v-text-field--is-booted v-select v-autocomplete primary--text"
  >
    <div class="v-input__control">
      <div class="v-input__slot">
        <div class="v-text-field__slot">
          <v-layout class="c-advanced-search__groups-wrapper gap-1" align-center wrap>
            <advanced-search-group
              v-for="(rule, index) in rules"
              v-model="rules[index]"
              :key="rule.key"
              :attributes="items"
              :active-key="activeKey"
              :union="index % 2 === 1"
              :first="index === 0"
              :allow-or="allowOr"
              disabled
              @input="update($event, index)"
              @click:chip="makeActive"
              @focusout="resetActive"
              @next="next($event, index)"
              @remove="remove(index)"
            />
          </v-layout>
        </div>
        <div class="v-input__append-inner">
          <v-menu bottom>
            <template #activator="{ on }">
              <c-action-btn
                :tooltip="$t('common.search')"
                icon="history"
                v-on="on"
              />
            </template>
          </v-menu>
        </div>
      </div>
    </div>
    <div>
      <c-action-btn
        :tooltip="$t('common.search')"
        icon="search"
        @click="submit"
      />
      <c-action-btn
        :tooltip="$t('common.clearSearch')"
        icon="clear"
        @click="clear"
      />
    </div>
  </div>
</template>

<script>
import AdvancedSearchGroup from '@/components/common/search/partials/new/advanced-search-group.vue';

export default {
  components: { AdvancedSearchGroup }

};
</script>

<style lang="scss" scoped>
.c-advanced-search { // TODO: remove new
  --v-chip-gap: 4px;
  --input-min-inline-size: 20ch;

  &__groups-wrapper > * {
    flex: 0 1 auto;
  }

  &::v-deep {
    input {
      flex: 0 1 auto;
      field-sizing: content;
      min-inline-size: var(--input-min-inline-size);
    }

    .layout {
      padding: var(--v-chip-gap) 0;
      gap: var(--v-chip-gap);
    }

    .v-chip {
      padding: 0 8px;
      margin: 0;

      &:has(> .v-chip__content > .v-chip) {
        padding: 0 6px !important;

        button {
          margin: 0 -2px 0 0 !important;
        }
      }

      &__content {
        gap: var(--v-chip-gap);
      }

      .v-chip {
        height: 24px;

        &.theme--light {
          background: var(--v-application-background-base, #FFFFFF);
        }

        &.theme--dark {
          background: var(--v-application-background-base, #121212);
        }
      }
    }

    button {
      margin-left: 4px !important;
    }
  }
}
</style>
