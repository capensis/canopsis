<template>
  <div class="footer">
    <a :href="canopsisDocumentation" target="_blank">
      {{ $t('login.documentation') }}
    </a>
    <a :href="canopsisWebsite" target="_blank">
      {{ $t('login.website') }}
    </a>
    <a :href="canopsisForum" target="_blank">
      {{ $t('login.forum') }}
    </a>
    <a :href="changeLogHref" class="version" target="_blank">
      {{ version }}
    </a>
  </div>
</template>

<script setup>
import { computed } from 'vue';

import { CANOPSIS_DOCUMENTATION, CANOPSIS_WEBSITE, CANOPSIS_FORUM } from '@/constants';

import { useInfo } from '@/hooks/store/modules/info';

const { version } = useInfo();

const changeLogHref = computed(() => {
  if (version.value) {
    return `${CANOPSIS_DOCUMENTATION}/notes-de-version/${version.value}/`;
  }

  return CANOPSIS_DOCUMENTATION;
});

const canopsisDocumentation = computed(() => CANOPSIS_DOCUMENTATION);
const canopsisWebsite = computed(() => CANOPSIS_WEBSITE);
const canopsisForum = computed(() => CANOPSIS_FORUM);
</script>

<style lang="scss">
.footer {
  grid-area: footer;
  position: relative;
  background-color: color-mix(in srgb, rgb(var(--v-theme-secondary)) 85%, black);
  color: white;
  min-height: 5em;
  margin-top: auto;
  display: flex;
  justify-content: center;
  align-items: center;

  a {
    color: inherit;
    text-decoration: none;
    padding: 0 2em;
  }

  .version {
    line-height: 5em;
    position: absolute;
    right: 0.5em;
    bottom: 0;
    text-decoration: underline;
    font-weight: bold;

    color: inherit;
    text-decoration: none;
  }
}
</style>
