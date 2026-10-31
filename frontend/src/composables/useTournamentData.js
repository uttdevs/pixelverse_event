/**
 * ============================================================================
 * COMPOSABLE: useTournamentData.js
 * ============================================================================
 * Centraliza la información del catálogo del torneo:
 * - Arenas de juego oficiales (con tipo de icono semántico).
 * - Categorías para filtrado.
 * - Podio de premios ($50,000 USD).
 * - Cronograma oficial y reglamento de torneo.
 */

/**
 * ============================================================================
 * COMPOSABLE: useTournamentData.js
 * ============================================================================
 * Centraliza la información del catálogo del torneo:
 * - Arenas de juego oficiales sincronizadas con /api/juegos.
 * - Mapa de estilos visuales gamer (iconos, gradientes, efectos).
 * - Cronograma oficial y reglamento de torneo.
 */

import { ref } from 'vue'

const ARENA_STYLES = {
  1: { // King of Fighters
    iconType: 'flame',
    accentColor: 'cyber-purple',
    glowClass: 'hover:box-glow-purple',
    bgGradient: 'from-purple-950 via-indigo-900 to-cyber-dark',
    buttonText: 'UNIRSE A LA ARENA'
  },
  2: { // Street Fighter
    iconType: 'bomb',
    accentColor: 'cyber-blue',
    glowClass: 'hover:box-glow-blue',
    bgGradient: 'from-blue-950 via-slate-900 to-cyber-dark',
    buttonText: 'UNIRSE A LA ARENA'
  },
  3: { // FIFA 2026
    iconType: 'crosshair',
    accentColor: 'cyber-blue',
    glowClass: 'hover:box-glow-blue',
    bgGradient: 'from-indigo-950 via-purple-900 to-cyber-dark',
    buttonText: 'UNIRSE A LA ARENA'
  }
}

const DEFAULT_STYLE = {
  iconType: 'crosshair',
  accentColor: 'cyber-blue',
  glowClass: 'hover:box-glow-blue',
  bgGradient: 'from-indigo-950 via-purple-900 to-cyber-dark',
  buttonText: 'UNIRSE A LA ARENA'
}

const getStyleForGame = (id, name) => {
  if (ARENA_STYLES[id]) return ARENA_STYLES[id]
  const lower = String(name || '').toLowerCase()
  if (lower.includes('fighter') && !lower.includes('street')) return ARENA_STYLES[1]
  if (lower.includes('street')) return ARENA_STYLES[2]
  if (lower.includes('fifa')) return ARENA_STYLES[3]
  return DEFAULT_STYLE
}

const defaultArenas = [
  {
    id_juego: 1,
    id: 'kof',
    name: 'THE KING OF FIGHTERS',
    formatBadge: 'Inscripción : $ 20',
    prize: '',
    description: 'Solo los más fuertes sobreviven a la arena, perfecciona tus combos y demuestra tu dominio en combates. ¡Entra en combate y demuestra quién merece el trono!',
    registeredParticipants: 0,
    ...ARENA_STYLES[1]
  },
  {
    id_juego: 2,
    id: 'sf',
    name: 'STREET FIGHTER',
    formatBadge: 'Inscripción : $ 20',
    prize: '',
    description: '¡Un rival, una batalla y una leyenda! Demuestra tu habilidad frente a los mejores combatientes, cada ronda es una prueba de técnica y concentración.',
    registeredParticipants: 0,
    ...ARENA_STYLES[2]
  },
  {
    id_juego: 3,
    id: 'fifa',
    name: 'FIFA 2026',
    formatBadge: 'Inscripción : $ 30',
    prize: '',
    description: 'El fútbol se juega con talento y se gana con estrategia. ¡Controla el balón y domina la cancha!',
    registeredParticipants: 0,
    ...ARENA_STYLES[3]
  }
]

export function useTournamentData() {
  const gameArenas = ref([...defaultArenas])
  const isLoadingArenas = ref(false)
  const arenasError = ref(null)

  const loadArenas = async () => {
    isLoadingArenas.value = true
    arenasError.value = null
    try {
      const res = await fetch('/api/juegos')
      if (res.ok) {
        const data = await res.json()
        if (Array.isArray(data) && data.length > 0) {
          gameArenas.value = data.map((j) => {
            const style = getStyleForGame(j.id_juego, j.nombre)
            const precio = j.costo_inscripcion != null ? Math.round(j.costo_inscripcion) : 20
            return {
              id_juego: j.id_juego,
              id: `juego-${j.id_juego}`,
              name: j.nombre,
              description: j.descripcion,
              formatBadge: `Inscripción : $ ${precio}`,
              prize: '',
              registeredParticipants: j.total_inscritos || 0,
              cuposMaximos: j.cupos_maximos,
              ...style
            }
          })
        }
      }
    } catch (err) {
      console.warn('No se pudo conectar a /api/juegos, usando fallback:', err)
      arenasError.value = err
    } finally {
      isLoadingArenas.value = false
    }
  }

  const scheduleStages = [
    {
      day: '20',
      month: 'NOVIEMBRE',
      status: '● FASE ACTIVA',
      statusColor: 'text-emerald-400 font-bold',
      borderHighlight: 'border-l-4 border-l-cyber-blue',
      badgeBg: 'bg-cyber-blue/10 border-cyber-blue/30 text-cyber-blue',
      title: 'CIERRE DE INSCRIPCIÓN',
      description: ''
  
    },
    {
      day: '27',
      month: 'NOVIEMBRE',
      status: 'FASE DE JUEGOS',
      statusColor: 'text-cyber-blue',
      badgeBg: 'bg-white/5 border-white/10 text-white',
      title: 'TORNEO DE VIDEOJUEGOS',
      description: ' 11 : 00 AM  a 1: 00 PM '
    },
    {
      day: '27',
      month: 'NOVIEMBRE',
      status: 'PREMIACIÓN',
      statusColor: 'text-cyber-purple font-bold',
      badgeBg: 'bg-cyber-purple/15 border-cyber-purple/40 text-purple-300',
      title: 'PREMIACIÓN Y CIERRE DEL TORNEO',
      description: '01:00 pm'
    }
  ]

  // Reglas del Torneo para RulesView
  const tournamentRules = [
    {
      section: '1. Reglas generales',
      items: [
        'Todos los participantes deberán respetar a sus rivales, organizadores y espectadores.',
        'La inscripción deberá realizarse antes del inicio de la competencia.',
        'Cada jugador deberá presentarse puntualmente en el horario asignado.',
        'Cualquier conducta antideportiva podrá ocasionar una advertencia o descalificación..',
        'Las decisiones de los organizadores deberán respetarse durante el desarrollo del torneo.',
        'Al inscribirse, los participantes aceptan las reglas establecidas.',
        'Queda prohibido utilizar trampas, modificaciones no autorizadas o cualquier método que otorgue una ventaja injusta.',
        '',
      ]
    },
   
    
    {
      section: 'Formato Competitivo',
      items: [
        'Los enfrentamientos se determinarán mediante el método establecido por la organización.',
        'Los participantes deberán consultar los horarios.',
        'Los resultados serán registrados para mantener el control de la competencia.'
      ]
    },
    {
      section: 'Faltas y sanciones',
      items: [
       
        ' Insultar, amenazar o acosar a otros participantes.',
        'Utilizar trampas o herramientas no autorizadas.',
'Interrumpir deliberadamente los enfrentamientos.',
' Suplantar a otro participante.',
' No presentarse a la hora asignada sin justificación.',
' Incumplir reiteradamente las indicaciones de la organización.'
      ]
    }
  
    
  ]

  return {
    gameArenas,
    isLoadingArenas,
    arenasError,
    loadArenas,
    scheduleStages,
    tournamentRules
  }
}
