package Publishing

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-co-op/gocron/v2"
	"os"
	"strconv"
	"time"
)

// func StarCronjob(ctx context.Context, serv Service, elasticServ ServiceElastic, mc *memcache.Client) error {
func StarCronjob(ctx context.Context, serv Service, elasticServ ServiceElastic) error {
	//Get enviroments variables
	cron_duration_job := os.Getenv("CRON_DURATION_JOB")
	// Fecha límite (solo fecha, sin hora)
	cron_job_deadline := os.Getenv("CRON_JOB_DEAD_LINE")

	MemDelta := os.Getenv("DELTA_INIT")

	durationJob, _ := strconv.Atoi(cron_duration_job)

	//Setting Timer
	// Obtener la fecha actual sin hora
	now := time.Now().Truncate(24 * time.Hour)

	// Parsear la fecha límite, el layout 2001-01-01
	deadline, err := time.Parse("2006-01-02", cron_job_deadline)
	if err != nil {
		fmt.Println("Error al parsear la fecha límite:", err)
		return errors.New("Error al parsear la fecha límite:")
	}

	// Calcular la diferencia de tiempo entre la fecha actual y la fecha límite
	diff := deadline.Sub(now)

	// Si la fecha actual es posterior a la fecha límite, detener la ejecución inmediatamente
	if diff <= 0 {
		fmt.Println("Fecha límite excedida, deteniendo la ejecución")
		return errors.New("Fecha límite excedida, deteniendo la ejecución")
	}

	// Crear un canal para el temporizador
	timer := time.NewTimer(diff)

	//SETTING CRON
	s, err := gocron.NewScheduler()
	if err != nil {
		fmt.Println("Error creating the scheduler")
	}

	j, err := s.NewJob(
		gocron.DurationJob(
			time.Duration(durationJob)*time.Minute,
		),
		gocron.NewTask(
			func(f string, h string) {

				/*				delta, errm := mc.Get("delta")
								fmt.Println(string(delta.Value))
								if errm != nil {
									// Actualiza el delta de la memoria cache
									nowTime := time.Now()
									// Formatear la fecha y hora actual en el formato YYYYDDMMHHMMSS
									d := nowTime.Format("20060102150405")
									err := SearchChanges(ctx, serv, elasticServ, d)
									if err != nil {
										fmt.Println("Error ", err)
									}
								} else {
									err := SearchChanges(ctx, serv, elasticServ, string(delta.Value))
									if err != nil {
										fmt.Println("Error ", err)
									}
								}

								// Actualiza el delta de la memoria cache
								nowTime := time.Now()
								// Formatear la fecha y hora actual en el formato YYYYDDMMHHMMSS
								formattedDelta := nowTime.Format("20060102150405")
								err = mc.Replace(&memcache.Item{Key: "delta", Value: []byte(formattedDelta)})
								if err != nil {
									log.Fatal(err)
								}
								delta, _ = mc.Get("delta")
								fmt.Println(string(delta.Value))
				*/
				err := SearchChanges(ctx, serv, elasticServ, MemDelta)
				if err != nil {
					fmt.Println("Error ", err)
				} else {
					// Actualiza el delta de la memoria cache
					nowTime := time.Now()
					// Formatear la fecha y hora actual en el formato YYYYDDMMHHMMSS
					MemDelta = nowTime.Format("20060102150405")

				}

			},
			"",
			"",
		),
	)

	if err != nil {
		fmt.Println("Error creating the new job: ", err)
	}
	// each job has a unique id
	fmt.Println(j.ID())

	// start the scheduler
	s.Start()

	// block until you are ready to shut down
	select {
	case <-timer.C:
		fmt.Println("Fecha límite del job alcanzada, deteniendo la ejecución del cron job")
	}

	// when you're done, shut it down
	err = s.Shutdown()
	if err != nil {
		fmt.Println("Error at shutdown")
	}

	return nil
}
