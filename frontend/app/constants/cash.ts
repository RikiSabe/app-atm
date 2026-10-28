export type CategoryId = 'constante' | 'variable'
export type MovementType = 'entrada' | 'salida'
export type PaymentMethod = 'QR' | 'Efectivo'
export type CashStatus = 'opening' | 'open' | 'closing' | 'closed'

export interface CashCategory {
  id: CategoryId
  name: string
}

export interface CashMovement {
  id: number
  category: CategoryId
  type: MovementType
  paymentMethod: PaymentMethod
  amount: number
  concept: string
  time: string
  inventoryItemId?: string
  inventoryQuantity?: number
}

export interface CashMovementDetail extends CashMovement {
  number: number
  balanceAfter: number
}

export interface PaymentSummary {
  count: number
  entries: number
  exits: number
  net: number
}

export const CATEGORIES: CashCategory[] = [
  { id: 'constante', name: 'Activo constante' },
  { id: 'variable', name: 'Activo variable' }
]

export const PAYMENT_METHODS: PaymentMethod[] = ['QR', 'Efectivo']

export function localDateValue() {
  const now = new Date()
  const month = String(now.getMonth() + 1).padStart(2, '0')
  const day = String(now.getDate()).padStart(2, '0')
  return `${now.getFullYear()}-${month}-${day}`
}

export function formatAmount(amount: number) {
  return new Intl.NumberFormat('es', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  }).format(amount)
}

export function formatSignedAmount(amount: number) {
  return amount < 0
    ? `−$${formatAmount(Math.abs(amount))}`
    : `+$${formatAmount(amount)}`
}

export function formatTime(time: string) {
  return new Date(time).toLocaleTimeString('es', {
    hour: '2-digit',
    minute: '2-digit'
  })
}

export function formatDateTime(time: string) {
  return new Date(time).toLocaleString('es', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

export function categoryName(categoryId: CategoryId) {
  return CATEGORIES.find(category => category.id === categoryId)?.name ?? categoryId
}

export function movementDetails(
  movements: CashMovement[],
  openingBalances: Record<CategoryId, number>
): CashMovementDetail[] {
  const balances = { ...openingBalances }

  return [...movements]
    .sort((first, second) => Date.parse(first.time) - Date.parse(second.time))
    .map((movement, index) => {
      balances[movement.category] += movement.type === 'entrada'
        ? movement.amount
        : -movement.amount

      return {
        ...movement,
        number: index + 1,
        balanceAfter: balances[movement.category]
      }
    })
}

export function paymentSummary(
  movements: CashMovement[],
  paymentMethod: PaymentMethod
): PaymentSummary {
  const methodMovements = movements.filter(movement => movement.paymentMethod === paymentMethod)
  const entries = methodMovements
    .filter(movement => movement.type === 'entrada')
    .reduce((total, movement) => total + movement.amount, 0)
  const exits = methodMovements
    .filter(movement => movement.type === 'salida')
    .reduce((total, movement) => total + movement.amount, 0)

  return {
    count: methodMovements.length,
    entries,
    exits,
    net: entries - exits
  }
}
