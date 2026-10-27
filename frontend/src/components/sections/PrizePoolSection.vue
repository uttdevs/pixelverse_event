<script setup>
/**
 * ============================================================================
 * COMPONENTE: PrizePoolSection.vue (Bolsa de Premios y Podio Gamer)
 * ============================================================================
 * Emplea iconos Lucide (Medal, Trophy, Award, Check, Star)
 * y directiva v-for con enlace dinámico de estilos.
 */

import { Medal, Trophy, Award, Check, Star } from 'lucide-vue-next'
import { useTournamentData } from '../../composables/useTournamentData'
import { useAudioSfx } from '../../composables/useAudioSfx'

const { prizeTiers } = useTournamentData()
const { playHoverSound } = useAudioSfx()
</script>


