<script setup>
/**
 * ============================================================================
 * COMPONENTE: BracketSection.vue (Llaves de Playoffs y Simulador de Predicciones)
 * ============================================================================
 * Emplea iconos Lucide (Info, Shuffle, Crown, Zap) y reactividad declarativa pura.
 */

import { ref } from 'vue'
import { Info, Shuffle, Crown, Zap } from 'lucide-vue-next'
import { useAudioSfx } from '../../composables/useAudioSfx'

const { playHoverSound, playClickSound } = useAudioSfx()

const bracket = ref({
  semiA: {
    team1: { name: 'KRONOS VORTEX', score: 2 },
    team2: { name: 'SHADOW WOLVES', score: 1 },
    selected: 'KRONOS VORTEX'
  },
  semiB: {
    team1: { name: 'NEON SAMURAIS', score: 2 },
    team2: { name: 'ASTRA GAMING', score: 0 },
    selected: 'NEON SAMURAIS'
  },
  champion: 'KRONOS VORTEX'
})

const teamsPoolA = ['KRONOS VORTEX', 'SHADOW WOLVES', 'CYBER DRAGONS']
const teamsPoolB = ['NEON SAMURAIS', 'ASTRA GAMING', 'ZERO LAG SQUAD']

const advanceTeam = (semiKey, teamName) => {
  playClickSound()
  bracket.value[semiKey].selected = teamName
  bracket.value.champion = teamName
}

const selectChampion = (teamName) => {
  playClickSound()
  bracket.value.champion = teamName
}

const randomizeBracket = () => {
  playClickSound()
  const randomA = teamsPoolA[Math.floor(Math.random() * teamsPoolA.length)]
  const randomB = teamsPoolB[Math.floor(Math.random() * teamsPoolB.length)]

  bracket.value.semiA.selected = randomA
  bracket.value.semiB.selected = randomB
  bracket.value.champion = Math.random() > 0.5 ? randomA : randomB
}
</script>


