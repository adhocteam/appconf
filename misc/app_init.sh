#!/bin/sh
### BEGIN INIT INFO
# Provides:          app
# Required-Start:    $local_fs $network $named $time $syslog
# Required-Stop:     $local_fs $network $named $time $syslog
# Default-Start:     2 3 4 5
# Default-Stop:      0 1 6
# Description:       app
### END INIT INFO

APPCONF_BUCKET=$(cat /etc/app/appconf_bucket)
SERVICE=$(cat /etc/app/service)
ENVIRONMENT=$(cat /etc/app/environment)

APP_NAME="app"
USER="app"
GROUP="app"
NODE_ENV="NODE_ENV=production NODE_OPTIONS=--abort-on-uncaught-exception"
APP_DIR="/srv/$APP_NAME"
CFG_FILE="$APP_DIR/.env"
PID_DIR="/var/run/$APP_NAME"
PID_FILE="$PID_DIR/$APP_NAME.pid"
LOG_DIR="/var/log/app/$ENVIRONMENT/$SERVICE"
LOG_FILE="$LOG_DIR/init.log"
NODE_EXEC="/usr/bin/npm"
NODE_ARGS="start"
USAGE="Usage: $0 {start|stop|restart|status|app_conf} [--force]"
FORCE_OP=false

pid_file_exists() {
    [ -f "$PID_FILE" ]
}

get_pid() {
    echo "$(cat "$PID_FILE")"
}

is_running() {
    PID="$(get_pid)"
    [ -d /proc/$PID ]
}

start_it() {
    mkdir -p "$PID_DIR"
    chown -R $USER:$GROUP "$PID_DIR"
    touch $PID_FILE
    chown $USER:$GROUP "$PID_FILE"
    mkdir -p "$LOG_DIR"
    chown -R $USER:$GROUP "$LOG_DIR"

    echo "Starting $APP_NAME ..."
    su -s /bin/bash -c "cd $APP_DIR && $NODE_ENV $NODE_EXEC $NODE_ARGS 1>>$LOG_FILE 2>&1 & echo \$! > $PID_FILE" $USER
    echo "$APP_NAME started with pid $(get_pid)"
}

stop_process() {
    PID=$(get_pid)
    echo "Killing process $PID"
    pkill -P $PID
}

remove_pid_file() {
    echo "Removing pid file"
    rm -f "$PID_FILE"
}

start_app() {
    # If the application conf file does not exist
    # Pull the conf from S3 before starting the service
    if [[ ! -f "$CFG_FILE" ]]; then
        appconf
    fi

    if pid_file_exists
    then
        if is_running
        then
            PID=$(get_pid)
            echo "$APP_NAME already running with pid $PID"
            exit 1
        else
            echo "$APP_NAME stopped, but pid file exists"
            if [ $FORCE_OP = true ]
            then
                echo "Forcing start anyways"
                remove_pid_file
                start_it
            fi
        fi
    else
        start_it
    fi
}

stop_app() {
    if pid_file_exists
    then
        if is_running
        then
            echo "Stopping $APP_NAME ..."
            stop_process
            remove_pid_file
            echo "$APP_NAME stopped"
        else
            echo "$APP_NAME already stopped, but pid file exists"
            if [ $FORCE_OP = true ]
            then
                echo "Forcing stop anyways ..."
                remove_pid_file
                echo "$APP_NAME stopped"
            else
                exit 1
            fi
        fi
    else
        echo "$APP_NAME already stopped, pid file does not exist"
        exit 1
    fi
}

status_app() {
    if pid_file_exists
    then
        if is_running
        then
            PID=$(get_pid)
            echo "$APP_NAME running with pid $PID"
        else
            echo "$APP_NAME stopped, but pid file exists"
        fi
    else
        echo "$APP_NAME stopped"
    fi
}

appconf() {
    # Cleanup
    rm -f "$CFG_FILE.new"
    rm -rf /srv/service-conf

    # Pull application configuration from S3
    s3Result=$(aws s3 cp s3://$APPCONF_BUCKET/$SERVICE/$ENVIRONMENT/ /srv/service-conf --recursive 2>&1)
    if [[ ! $? -eq 0 ]]; then
        echo "$s3Result"
        logger "$s3Result"

        echo "Unable to download configuration from S3"
        logger "Unable to download configuration from S3"

        # Cleanup
        rm -rf /srv/service-conf

        exit 1
    fi

    # Setup the .new env file
    touch "$CFG_FILE.new"
    chmod 640 "$CFG_FILE.new"
    chown $USER:$GROUP "$CFG_FILE.new"

    # Output filename=filecontents to the new .env file
    for file in /srv/service-conf/*; do
        name=$(basename $file)
        printf "%s=%s\n" "$name" "$(<$file)" >> "$CFG_FILE.new"
    done

    # If the config file doesn't exist, create it and start the service
    if [[ ! -f "$CFG_FILE" ]]; then
        echo "Configuration file does not exist; creating it"
        logger "Configuration file does not exist; creating it"

        mv "$CFG_FILE.new" "$CFG_FILE"
    fi

    # If the config files are different, update and restart
    diffResult=$(diff "$CFG_FILE" "$CFG_FILE.new" 2>&1)
    if [[ $? -eq 1 ]]; then
        echo "App configuration file is different"
        logger "App configuration file is different"

        if [ $FORCE_OP = true ]; then
            echo "Replacing configuration file"
            logger "Replacing configuration file"
            mv "$CFG_FILE.new" "$CFG_FILE"

            if pid_file_exists && is_running; then
                echo "Restarting service"
                logger "Restarting service"

                stop_process

                remove_pid_file

                start_it
            fi
        fi
    else
        echo "App configuration file is OK"
    fi

    # Cleanup
    rm -f "$CFG_FILE.new"
    rm -rf /srv/service-conf
}

case "$2" in
    --force)
        FORCE_OP=true
    ;;

    "")
    ;;

    *)
        echo $USAGE
        exit 1
    ;;
esac

case "$1" in
    start)
        start_app
    ;;

    stop)
        stop_app
    ;;

    restart)
        stop_app
        start_app
    ;;

    status)
        status_app
    ;;

    appconf)
        appconf
    ;;

    *)
        echo $USAGE
        exit 1
    ;;
esac
