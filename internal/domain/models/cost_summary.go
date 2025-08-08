package models

/*
CostSummary — агрегатор для подсчёта общей стоимости подписок
за определённый период. Хранит:
- totalCost — общая сумма
- period — диапазон дат, за который ведётся расчёт
- subscriptions — список подписок, по которым идёт расчёт
*/
type CostSummary struct {
	totalCost     int
	period        DatePeriod
	subscriptions []Subscription
}

/** Создаёт новый объект для подсчёта с заданным периодом. */
func NewCostSummary(period DatePeriod) *CostSummary {
	return &CostSummary{
		period:        period,
		subscriptions: make([]Subscription, 0),
	}
}

/** Геттер/сеттер для общей суммы. */
func (cs *CostSummary) TotalCost() int {
	return cs.totalCost
}

func (cs *CostSummary) SetTotalCost(totalCost int) {
	cs.totalCost = totalCost
}

/** Геттер/сеттер для периода расчёта. */
func (cs *CostSummary) Period() DatePeriod {
	return cs.period
}

func (cs *CostSummary) SetPeriod(period DatePeriod) {
	cs.period = period
	cs.totalCost = 0 // сбрасываем сумму, так как период изменился
}

/** Геттер/сеттер для списка подписок. */
func (cs *CostSummary) Subscriptions() []Subscription {
	return cs.subscriptions
}

func (cs *CostSummary) SetSubscriptions(subscriptions []Subscription) {
	cs.subscriptions = subscriptions
}

/** Добавляет одну подписку в список. */
func (cs *CostSummary) AddSubscription(sub Subscription) {
	cs.subscriptions = append(cs.subscriptions, sub)
}

/*
*
Calculate — считает суммарную стоимость всех подписок
за указанный период, используя CalculateCostForPeriod каждой подписки.
Результат сохраняется в totalCost и возвращается.
*/
func (cs *CostSummary) Calculate() int {
	total := 0
	for _, sub := range cs.subscriptions {
		total += sub.CalculateCostForPeriod(cs.period.From(), cs.period.To())
	}
	cs.totalCost = total
	return total
}
