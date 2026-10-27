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

export function useTournamentData() {
  const gameArenas = [
    {
      id: 'valorant',
      name: 'FIFA 2026',
      category: 'fps',
      formatBadge: 'Inscripción : $ 30',
      prize: '',
      description:'El fútbol se juega con talento y se gana con estrategia.¡Controla el balón y domina la cancha! ',
     
      iconType: 'crosshair',
      accentColor: 'cyber-blue',
      glowClass: 'hover:box-glow-blue',
      bgGradient: 'from-indigo-950 via-purple-900 to-cyber-dark',
      buttonText: 'UNIRSE A LA ARENA'
    },
    {
      id: 'lol',
      name: 'THE KING OF FIGHTERS ',
      category: 'moba',
      formatBadge: 'Inscripción : $ 20',
      prize: '',
      description: 'Solo los más fuertes sobreviven a la arena, perfecciona tus combos y demuestra tu dominio en combates.¡Entra en combate y demuestra quien merece el trono!',
     
      iconType: 'flame',
      accentColor: 'cyber-purple',
      glowClass: 'hover:box-glow-purple',
      bgGradient: 'from-purple-950 via-indigo-900 to-cyber-dark',
      buttonText: 'UNIRSE A LA ARENA'
    },
    {
      id: 'cs2',
      name: 'STREET FIGHTER',
      category: 'fps',
      formatBadge: 'Inscripción : $ 20',
      prize: '',
      description: '¡Un rival,una batalla y una leyenda!.Demuestra tu habilidad frente a los mejores combatientes,cada ronda es una prueba de técnica y concentración.',
      
      iconType: 'bomb',
      accentColor: 'cyber-blue',
      glowClass: 'hover:box-glow-blue',
      bgGradient: 'from-blue-950 via-slate-900 to-cyber-dark',
      buttonText: 'UNIRSE A LA ARENA'
    },
  
  ]

  const categories = [
    { id: 'all', label: 'TODOS (3)' },
   
  ]

  const prizeTiers = [
  
    
   
  ]

  const scheduleStages = [
    {
      day: '20',
      month: 'NOVIEMBRE',
      status: '● FASE ACTIVA',
      statusColor: 'text-emerald-400 font-bold',
      borderHighlight: 'border-l-4 border-l-cyber-blue',
      badgeBg: 'bg-cyber-blue/10 border-cyber-blue/30 text-cyber-blue',
      title: 'CIERRE DE INSCRIPCIÓN',
      description: 'Requisito para FIFA control inalámbrico para Xbox y Registro previo para participar.'
  
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
      description: 'Primer lugar de cada categoría ,el premio será monetario'
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
    categories,
    prizeTiers,
    scheduleStages,
    tournamentRules
  }
}
