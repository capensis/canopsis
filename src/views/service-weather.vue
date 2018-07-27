<template lang="pug">
  weather(:ids=["periode.moderne.check.resource.Les lumières.L'Encyclopédie","periode.historique.check.resource.Moyen âge.Serments de Strasbourg","periode.moderne.check.resource.Renaissance.L'École d'Athènes"])
</template>

<script>
import Weather from '@/components/other/service-weather/weather.vue';

export default {
  components: {
    Weather,
  },
};
</script>

