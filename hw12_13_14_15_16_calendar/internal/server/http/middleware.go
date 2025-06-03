package internalhttp

// Вопрос это большая зависимость на тип логера? но врядтли кто-то когда-то захочет поменять логгер на проекте?
//func New(log *slog.Logger) func(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		// TODO
//		next.ServeHTTP(w, r)
//	})
//}
